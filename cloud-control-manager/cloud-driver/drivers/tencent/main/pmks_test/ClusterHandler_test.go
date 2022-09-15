package pmks

import (
	"os"
	"testing"

	_ "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/drivers/tencent/main/pmks_test/env" // 위치 변경 하면 안됨. 환경설정 정보 읽기 전에 테스트 수행됨

	tdrv "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/drivers/tencent"
	idrv "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces"
	irs "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces/resources"
)

func getClusterHandler() (irs.ClusterHandler, error) {

	connectionInfo := idrv.ConnectionInfo{
		CredentialInfo: idrv.CredentialInfo{
			ClientId:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
		},
		RegionInfo: idrv.RegionInfo{
			Region: os.Getenv("REGION"),
			Zone:   os.Getenv("ZONE"),
		},
	}

	cloudDriver := new(tdrv.TencentDriver)
	cloudConnection, err := cloudDriver.ConnectCloud(connectionInfo)
	if err != nil {
		return nil, err
	}

	clusterHandler, err := cloudConnection.CreateClusterHandler()
	if err != nil {
		return nil, err
	}

	return clusterHandler, nil
}

func TestGetClusterHander(t *testing.T) {
	clusterHandler, err := getClusterHandler()
	if err != nil {
		t.Error(err)
	}

	println(clusterHandler)
}

func TestCreateClusterOnly(t *testing.T) {

	t.Log("클러스터 생성, 노드그룹은 생성안함")

	clusterHandler, err := getClusterHandler()
	if err != nil {
		t.Error(err)
	}

	clusterInfo := irs.ClusterInfo{
		IId: irs.IID{
			NameId:   "cluster-x1",
			SystemId: "",
		},
		Version: "1.22.5",
		Network: irs.NetworkInfo{
			VpcIID:    irs.IID{NameId: "", SystemId: "vpc-q1c6fr9e"},
			SubnetIID: []irs.IID{{NameId: "", SystemId: "subnet-rl79gxhv"}},
			//SecurityGroupIIDs: []irs.IID{{NameId: "", SystemId: "sg-46eef229"}}, // 설정 안됨
		},
	}

	cluster_, err := clusterHandler.CreateCluster(clusterInfo)
	if err != nil {
		t.Error(err)
	}

	t.Log(cluster_)
}

func TestListCluster(t *testing.T) {

	clusterHandler, err := getClusterHandler()
	if err != nil {
		t.Error(err)
	}

	clusters, err := clusterHandler.ListCluster()
	if err != nil {
		t.Error(err)
	}

	if len(clusters) == 0 {
		t.Error("No cluster found")
	}

	for _, cluster := range clusters {
		t.Log(cluster.IId.SystemId)
		println(cluster.IId.NameId, cluster.Status)
	}
}

func TestGetCluster(t *testing.T) {

	clusterHandler, err := getClusterHandler()
	if err != nil {
		t.Error(err)
	}

	clusters, err := clusterHandler.ListCluster()
	if err != nil {
		t.Error(err)
	}

	if len(clusters) == 0 {
		t.Error("No cluster found")
	}

	t.Log(clusters)

	for _, cluster := range clusters {
		cluster_, err := clusterHandler.GetCluster(cluster.IId)
		if err != nil {
			println(err.Error())
		}
		t.Log(cluster_)
	}
}

// 맨 마지막 테스트로 이동
// func TestDeleteCluster(t *testing.T) {
// }

func TestAddNodeGroup(t *testing.T) {
	clusterHandler, err := getClusterHandler()
	if err != nil {
		t.Error(err)
	}

	new_node_group := &irs.NodeGroupInfo{
		IId: irs.IID{NameId: "nodepoolx101", SystemId: ""},
		// image id can not be set, when creating nodepool
		// ImageIID:        irs.IID{NameId: "", SystemId: "img-pi0ii46r"}, // 이미지 id 선택 추가, img-pi0ii46r:ubuntu18.04
		VMSpecName:      "S3.MEDIUM2",
		RootDiskType:    "CLOUD_PREMIUM",
		RootDiskSize:    "50",
		KeyPairIID:      irs.IID{NameId: "kp1", SystemId: ""}, // 필수 옵션 아님, 대응되는 필드가 없음. 찾아봐야함.
		OnAutoScaling:   true,
		DesiredNodeSize: 1,
		MinNodeSize:     0,
		MaxNodeSize:     3,
	}

	clusters, _ := clusterHandler.ListCluster()
	for _, cluster := range clusters {
		t.Log(cluster)
		node_group, err := clusterHandler.AddNodeGroup(cluster.IId, *new_node_group)
		if err != nil {
			t.Error(err)
		}
		t.Log(node_group)
	}
}

func TestListNodeGroup(t *testing.T) {
	clusterHandler, err := getClusterHandler()
	if err != nil {
		t.Error(err)
	}

	clusters, _ := clusterHandler.ListCluster()
	for _, cluster := range clusters {
		node_groups, _ := clusterHandler.ListNodeGroup(cluster.IId)
		for _, node_group := range node_groups {
			t.Log(node_group.IId.NameId, node_group.IId.SystemId)
			t.Log(node_group)
		}
	}
}

func TestGetNodeGroup(t *testing.T) {
	clusterHandler, err := getClusterHandler()
	if err != nil {
		t.Error(err)
	}

	clusters, _ := clusterHandler.ListCluster()
	for _, cluster := range clusters {
		node_groups, _ := clusterHandler.ListNodeGroup(cluster.IId)
		for _, node_group := range node_groups {
			node_group_, err := clusterHandler.GetNodeGroup(cluster.IId, node_group.IId)
			if err != nil {
				t.Error(err)
			}
			t.Log(node_group_.IId.NameId, node_group_.IId.SystemId)
			t.Log(node_group_)
		}
	}

	node_group, err := clusterHandler.GetNodeGroup(irs.IID{NameId: "", SystemId: "cluster_id_not_exist"}, irs.IID{NameId: "", SystemId: "node_group_id_not_exist"})
	if err != nil {
		println(err.Error())
	}
	println(node_group.IId.NameId)
}

func TestSetNodeGroupAutoScaling(t *testing.T) {
	clusterHandler, err := getClusterHandler()
	if err != nil {
		t.Error(err)
	}

	clusters, _ := clusterHandler.ListCluster()
	for _, cluster := range clusters {
		node_groups, _ := clusterHandler.ListNodeGroup(cluster.IId)
		for _, node_group := range node_groups {
			node_group_, err := clusterHandler.GetNodeGroup(cluster.IId, node_group.IId)
			if err != nil {
				t.Error(err)
			}

			res, err := clusterHandler.SetNodeGroupAutoScaling(cluster.IId, node_group_.IId, false)
			if err != nil {
				t.Error(err)
			}
			println(res)

			res, err = clusterHandler.SetNodeGroupAutoScaling(cluster.IId, node_group_.IId, true)
			if err != nil {
				t.Error(err)
			}
			println(res)
		}
	}
}

func TestChangeNodeGroupScaling(t *testing.T) {
	clusterHandler, err := getClusterHandler()
	if err != nil {
		t.Error(err)
	}

	clusters, _ := clusterHandler.ListCluster()
	for _, cluster := range clusters {
		node_groups, _ := clusterHandler.ListNodeGroup(cluster.IId)
		for _, node_group := range node_groups {
			node_group_, err := clusterHandler.GetNodeGroup(cluster.IId, node_group.IId)
			if err != nil {
				t.Error(err)
			}

			res, err := clusterHandler.ChangeNodeGroupScaling(cluster.IId, node_group_.IId, 2, 0, 5)
			if err != nil {
				t.Error(err)
			}
			println(res.IId.NameId, res.IId.SystemId)

			res, err = clusterHandler.ChangeNodeGroupScaling(cluster.IId, node_group_.IId, 1, 0, 3)
			if err != nil {
				t.Error(err)
			}
			println(res.IId.NameId, res.IId.SystemId)
		}
	}
}

func TestRemoveNodeGroup(t *testing.T) {
	clusterHandler, err := getClusterHandler()
	if err != nil {
		t.Error(err)
	}

	clusters, _ := clusterHandler.ListCluster()
	for _, cluster := range clusters {
		node_groups, _ := clusterHandler.ListNodeGroup(cluster.IId)
		for _, node_group := range node_groups {
			res, _ := clusterHandler.RemoveNodeGroup(cluster.IId, node_group.IId)
			if err != nil {
				t.Error(err)
			}
			if res == false {
				t.Error("Failed to remove node group")
			}
		}
	}
}

func TestUpgradeCluster(t *testing.T) {
	clusterHandler, err := getClusterHandler()
	if err != nil {
		t.Error(err)
	}

	clusters, _ := clusterHandler.ListCluster()
	for _, cluster := range clusters {
		//res, err := clusterHandler.UpgradeCluster(cluster.IId, "1.22.5") //version := "1.22.5"
		res, err := clusterHandler.UpgradeCluster(cluster.IId, "1.20.6") //version := "1.22.5"
		if err != nil {
			t.Error(err)
		}
		t.Log(res)
	}
}

func TestDeleteCluster(t *testing.T) {

	clusterHandler, err := getClusterHandler()
	if err != nil {
		t.Error(err)
	}

	clusters, _ := clusterHandler.ListCluster()
	for _, cluster := range clusters {
		println(cluster.IId.NameId, cluster.IId.SystemId)
		result, err := clusterHandler.DeleteCluster(cluster.IId)
		if err != nil {
			t.Error(err)
		}
		t.Log(result)
	}

	// result, err := clusterHandler.DeleteCluster(irs.IID{NameId: "cluster_not_exist", SystemId: "cluster_id_not_exist"})
	// if err != nil {
	// 	println(err.Error())
	// }
	// println(result)
}
