package google

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/infracost/infracost/internal/resources/google"
	"github.com/infracost/infracost/internal/schema"

	"strings"
)

func getDataprocClusterRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "google_dataproc_cluster",
		RFunc: newDataprocCluster,
	}
}

func newDataprocCluster(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	region := d.Get("region").String()

	// In newDataprocCluster we parse the resource data coming from d
	// into the underlying DataprocCluster. We use schema.ResourceData lookup
	// to find data stored on attributes in the resource. This can be done as follows:
	sku := d.Get("sku_name").String()
	instanceType := strings.Split(sku, "_")[0]
	count := d.Get("terraformFieldThatHasNumberOfInstances").Int()

	usageType := "default"
	if !d.IsEmpty("usageType") {
		usageType = d.Get("usageType").String()
	}

	//containerMasterConfig := newContainerNodeConfig(d.Get("master_config.0"))
	//containerWorkerConfig := newContainerNodeConfig(d.Get("worker_config.0"))
	clusterConfigBucket := d.Get("cluster_config.0.staging_bucket")
	fmt.Printf("Bucket: \n%s\n\n", clusterConfigBucket)
	containerMasterConfig := newMasterConfig(d.Get("cluster_config.0.master_config.0"))
	containerMasterDiskConfig := newDiskConfig(d.Get("cluster_config.0.master_config.0.disk_config.0"))
	containerWorkerConfig := newMasterConfig(d.Get("cluster_config.0.worker_config.0"))
	containerWorkerDiskConfig := newDiskConfig(d.Get("cluster_config.0.worker_config.0.disk_config.0"))
	preemptibleWorkerConfig:= newMasterConfig(d.Get("cluster_config.0.preemptible_worker_config.0"))
	preemptibleWorkerDiskConfig:= newDiskConfig(d.Get("cluster_config.0.preemptible_worker_config.0.disk_config.0"))

	

	//fmt.Printf("MachineType: \n%s\n\n", r.MasterConfig.MachineType)
	//fmt.Printf("BootDiskType: \n%s\n\n", r.MasterDiskConfig.BootDiskType)
	

	r := &google.DataprocCluster{
		Address:       d.Address,
		Region:        region,
		InstanceCount: count,
		InstanceType:  instanceType,
		UsageType:     usageType,
		MasterConfig:   containerMasterConfig,
		MasterDiskConfig:   containerMasterDiskConfig,
		WorkerConfig:   containerWorkerConfig,
		WorkerDiskConfig:   containerWorkerDiskConfig,
		PreemptibleWorkerConfig: preemptibleWorkerConfig,
		PreemptibleWorkerDiskConfig: preemptibleWorkerDiskConfig,
	}
	r.PopulateUsage(u)

	return r.BuildResource()
}

func newMasterConfig(d gjson.Result) *google.MasterConfig {
	machineType := "e2-medium"
	fmt.Printf("MachineType: \n%s\n\n", d.Get("machine_type").String())
	if d.Get("machine_type").Exists() {
		machineType = d.Get("machine_type").String()
	}

	purchaseOption := "on_demand"
	if d.Get("preemptible").Bool() {
		purchaseOption = "preemptible"
	}

	numInstances := d.Get("num_instances").Int()

	guestAccelerators := collectComputeGuestAccelerators(d.Get("guest_accelerator"))

	return &google.MasterConfig{
		MachineType:       machineType,
		PurchaseOption:    purchaseOption,
		NumInstances:      numInstances,
		GuestAccelerators: guestAccelerators,
	}
}

func newDiskConfig(d gjson.Result) *google.DiskConfig {

	diskType := "pd-standard"
	if d.Get("boot_disk_type").Exists() {
		diskType = d.Get("boot_disk_type").String()
	}

	diskSize := int64(100)
	if d.Get("boot_disk_size_gb").Exists() {
		diskSize = d.Get("boot_disk_size_gb").Int()
	}

	numLocalSSDs := d.Get("num_local_ssds").Int()

	return &google.DiskConfig{
		BootDiskType:          diskType,
		BootDiskSize:          float64(diskSize),
		NumLocalSSDs:     numLocalSSDs,
	}
}

