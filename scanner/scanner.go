package scanner

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/opsee/basic/schema"
	opsee_aws "github.com/opsee/basic/schema/aws"
	"net"
	"sort"
)

func ScanRegion(region string, session *session.Session) (*schema.Region, error) {
	ec2Client := ec2.New(session)

	var (
		nextToken      *string
		vpcIPs         = make(map[string][]string)
		vpcRouteTables = make(map[string][]*ec2.RouteTable)
		vpcGateways    = make(map[string][]*ec2.InternetGateway)
		vpcCounts      = make(map[string]int)
		subnetCounts   = make(map[string]int)
	)

	for {
		instancesOutput, err := ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
			MaxResults: aws.Int64(100),
			NextToken:  nextToken,
		})

		if err != nil {
			return nil, err
		}

		nextToken = instancesOutput.NextToken
		for _, res := range instancesOutput.Reservations {
			for _, instance := range res.Instances {
				if aws.StringValue(instance.State.Name) != ec2.InstanceStateNameTerminated {
					vips, ok := vpcIPs[aws.StringValue(instance.VpcId)]
					if !ok {
						vips = make([]string, 0)
					}

					vpcIPs[aws.StringValue(instance.VpcId)] = append(vips, aws.StringValue(instance.PrivateIpAddress))
					vpcCounts[aws.StringValue(instance.VpcId)]++
					subnetCounts[aws.StringValue(instance.SubnetId)]++
				}
			}
		}

		if nextToken == nil {
			break
		}
	}

	internetGatewaysOutput, err := ec2Client.DescribeInternetGateways(nil)
	if err != nil {
		return nil, err
	}

	for _, igw := range internetGatewaysOutput.InternetGateways {
		for _, igwatt := range igw.Attachments {
			st := aws.StringValue(igwatt.State)
			if st == attachmentStatusAvailable || st == ec2.AttachmentStatusAttached {
				vigw, ok := vpcGateways[aws.StringValue(igwatt.VpcId)]
				if !ok {
					vigw = make([]*ec2.InternetGateway, 0)
				}

				vpcGateways[aws.StringValue(igwatt.VpcId)] = append(vigw, igw)
			}
		}
	}

	routeTablesOutput, err := ec2Client.DescribeRouteTables(nil)
	if err != nil {
		return nil, err
	}

	for _, rt := range routeTablesOutput.RouteTables {
		vrts, ok := vpcRouteTables[aws.StringValue(rt.VpcId)]
		if !ok {
			vrts = make([]*ec2.RouteTable, 0)
		}

		vpcRouteTables[aws.StringValue(rt.VpcId)] = append(vrts, rt)
	}

	vpcOutput, err := ec2Client.DescribeVpcs(nil)
	if err != nil {
		return nil, err
	}

	vpcs := make([]*schema.Vpc, len(vpcOutput.Vpcs))
	for vi, v := range vpcOutput.Vpcs {
		vpc := &schema.Vpc{}
		opsee_aws.CopyInto(vpc, v)

		vpc.InstanceCount = int32(vpcCounts[vpc.VpcId])
		vpcs[vi] = vpc
	}

	subnetOutput, err := ec2Client.DescribeSubnets(nil)
	if err != nil {
		return nil, err
	}

	subnets := make([]*schema.Subnet, len(subnetOutput.Subnets))
	for si, s := range subnetOutput.Subnets {
		subnet := &schema.Subnet{}
		opsee_aws.CopyInto(subnet, s)

		subnet.InstanceCount = int32(subnetCounts[subnet.SubnetId])

		routing, err := determineRouting(s, vpcRouteTables[aws.StringValue(s.VpcId)], vpcIPs[aws.StringValue(s.VpcId)], vpcGateways[aws.StringValue(s.VpcId)])
		if err != nil {
			return nil, err
		}

		subnet.Routing = routing
		subnets[si] = subnet
	}

	// now sort them in the order we want to select for bastion install
	sort.Sort(schema.SubnetsByPreference(subnets))

	accountOutput, err := ec2Client.DescribeAccountAttributes(&ec2.DescribeAccountAttributesInput{
		AttributeNames: []*string{
			aws.String("supported-platforms"),
		},
	})
	if err != nil {
		return nil, err
	}

	var supportedPlatforms []string
	for _, a := range accountOutput.AccountAttributes {
		if aws.StringValue(a.AttributeName) == "supported-platforms" {
			for _, v := range a.AttributeValues {
				supportedPlatforms = append(supportedPlatforms, aws.StringValue(v.AttributeValue))
			}
		}
	}

	return &schema.Region{
		Region:             region,
		SupportedPlatforms: supportedPlatforms,
		Vpcs:               vpcs,
		Subnets:            subnets,
	}, nil
}

func determineRouting(subnet *ec2.Subnet, routeTables []*ec2.RouteTable, instanceIPs []string, gateways []*ec2.InternetGateway) (string, error) {
	var (
		associatedTable *ec2.RouteTable
		mainTable       *ec2.RouteTable
		internetRoute   *ec2.Route
		cidrs           = make([]string, 0)
		routeToIPs      = make(map[string]bool)
	)

	for _, rt := range routeTables {
		for _, asso := range rt.Associations {
			if aws.BoolValue(asso.Main) {
				mainTable = rt
			}

			if asso.SubnetId != nil && aws.StringValue(asso.SubnetId) == aws.StringValue(subnet.SubnetId) {
				associatedTable = rt
			}
		}
	}

	if associatedTable == nil {
		associatedTable = mainTable
	}

	// there is no main table or associated tables, so who knows what is going on
	if associatedTable == nil {
		return schema.RoutingStatePrivate, nil
	}

	for _, route := range associatedTable.Routes {
		if aws.StringValue(route.State) != ec2.RouteStateActive {
			continue
		}

		if route.DestinationCidrBlock == nil {
			continue
		}

		if aws.StringValue(route.DestinationCidrBlock) == theInternet {
			internetRoute = route
			continue
		}

		cidrs = append(cidrs, aws.StringValue(route.DestinationCidrBlock))
	}

	// no route to 0.0.0.0/0, so we're private. no need to check if
	// we have a route to other instances
	if internetRoute == nil {
		return schema.RoutingStatePrivate, nil
	}

	// TODO: scan network ACLs here? YES. going to push this out first tho

	// verify that we can reach all of the instance ips? i say yes
	for _, cid := range cidrs {
		_, network, err := net.ParseCIDR(cid)
		if err != nil {
			return "", err
		}

		for _, ip := range instanceIPs {
			if network.Contains(net.ParseIP(ip)) {
				routeToIPs[ip] = true
			}
		}
	}

	// we must not be able to reach all the instance ips,
	// so we're going to mark this subnet as occluded
	if len(routeToIPs) < len(instanceIPs) {
		return schema.RoutingStateOccluded, nil
	}

	// nat. pretty straightforward i guess,
	// going to punt on verifying that the nat instance itself has a route?
	if internetRoute.InstanceId != nil {
		return schema.RoutingStateNAT, nil
	}

	// public _or_ routing through a customer gateway, in which case we'll consider it NAT,
	// since the only case we'll need a public ip is if we're going through an internet gateway
	if internetRoute.GatewayId != nil {
		for _, igw := range gateways {
			if aws.StringValue(internetRoute.GatewayId) == aws.StringValue(igw.InternetGatewayId) {
				return schema.RoutingStatePublic, nil
			}
		}

		return schema.RoutingStateGateway, nil
	}

	// we're in a weird state. the routing table is possibly attached to an eni, but not an instance
	return schema.RoutingStateOccluded, nil
}
