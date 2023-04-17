package google

import (
	"fmt"
	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

// DataprocCluster struct represents <TODO: cloud service short description>.
//
// <TODO: Add any important information about the resource and links to the
// pricing pages or documentation that might be useful to developers in the future, e.g:>
//
// Resource information: https://cloud.google.com/<PATH/TO/RESOURCE>/
// Pricing information: https://cloud.google.com/<PATH/TO/PRICING>/
type DataprocCluster struct {
	// Add DataprocCluster struct properties below that correspond to your resource
	// and are directly used in cost calculation.
	//
	// These properties are configured directly from attributes parsed in the IaC provider.
	// See your resource file: internal/providers/terraform/google/dataproc_cluster.go
	// for an example in how this is achieved.
	//
	// Below there are a few examples of common properties that are often part of resources:
	// Address is the unique name of the resource in the IaC language. It is required.
	Address string
	// Region is the google region the DataprocCluster is provisioned within. It is required.
	Region string

	// The properties below are examples to show how they can be used with cost components.
	InstanceCount int64
	InstanceType  string
	UsageType     string
	MasterConfig   *MasterConfig
	WorkerConfig   *MasterConfig
	MasterDiskConfig *DiskConfig
	WorkerDiskConfig *DiskConfig
	PreemptibleWorkerConfig *MasterConfig
	PreemptibleWorkerDiskConfig *DiskConfig

	// Add the usage parameters for DataprocCluster below.
	//
	// A usage parameter is defined simply as a property on the main resource struct.
	// But it needs to have `infracost_usage` so that the `PopulateUsage` method can
	// extract the correct value from the usage file.
	//
	// Below is an example usage parameter MonthlyDataProcessedGB. Feel free to delete this if it's not needed in DataprocCluster.
	//
	// This property would work if you have a usage parameter defined in your usage file as such:
	//
	//   	google_dataproc_cluster.dataproc_cluster:
	//      monthly_data_processed_gb: 200.50
	//
	// It should be defined as a parameter on this struct like the one below.
	// Note the `infracost_usage` tag matches the name of the property in the usage file.

	// "usage" args
	MonthlyDataProcessedGB *float64 `infracost_usage:"monthly_data_processed_gb"`
}

// DataprocClusterUsageSchema defines a list which represents the usage schema of DataprocCluster.
// If DataprocCluster has no usage schema it's safe to delete this type.
var DataprocClusterUsageSchema = []*schema.UsageItem{
	// e.g. if you follow the example given in DataprocCluster struct you would need to add a list item with the following:
	{Key: "monthly_data_processed_gb", DefaultValue: 0, ValueType: schema.Float64},
	// Replace the above with all the usage items you need for DataprocCluster.
}

// PopulateUsage parses the u schema.UsageData into the DataprocCluster.
// It uses the `infracost_usage` struct tags to populate data into the DataprocCluster.
func (r *DataprocCluster) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid DataprocCluster struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *DataprocCluster) BuildResource() *schema.Resource {
	// List of schema.CostComponent items is how Infracost builds price outputs.
	// Below are a few examples based on the dummy parameters defined in the struct above.
	//
	// These are examples and are safe to remove.
	//costComponents := []*schema.CostComponent{
		// Below is an example cost component using properties that we've parsed
		// from the IaC into the DataprocCluster.InstanceCount & DataprocCluster.InstanceType.
	//	r.instanceCostComponent(),
		// Below is an example of a cost component built with the parsed usage property.
		// Note the r.MonthlyDataProcessedGB field passed to hourly quantity.
	//	r.dataProcessedCostComponent(),
	//}

	poolSize := int64(1)

	fmt.Printf("Region: \n%s\n\n", r.Region)
	fmt.Printf("MachineType: \n%s\n\n", r.MasterConfig.MachineType)
	fmt.Printf("BootDiskType: \n%s\n\n", r.MasterDiskConfig.BootDiskType)

	//Master Node Config
	costComponents := []*schema.CostComponent{
		//computeCostComponent(r.Region, r.MasterConfig.MachineType, r.MasterConfig.PurchaseOption, r.MasterConfig.NumInstances, nil),
		computeCostDataProcComponent(r.Region, r.MasterConfig.MachineType, r.MasterConfig.PurchaseOption, r.MasterConfig.NumInstances, nil,"Master Node"),
		computeDiskCostComponent(r.Region, r.MasterDiskConfig.BootDiskType, r.MasterDiskConfig.BootDiskSize, poolSize),
		//computeDiskCostComponent(r.Region, "pd-standard", 200, poolSize),
	}

	if r.MasterDiskConfig.NumLocalSSDs > 0 {
		costComponents = append(costComponents, scratchDiskCostComponent(r.Region, r.MasterConfig.PurchaseOption, int(r.MasterDiskConfig.NumLocalSSDs)))
	}

	//Worker Node Config
	costComponents = append(costComponents, computeCostDataProcComponent(r.Region, r.WorkerConfig.MachineType, r.WorkerConfig.PurchaseOption, 1, nil,"Worker Node"))

	costComponents = append(costComponents, computeDiskCostComponent(r.Region, r.WorkerDiskConfig.BootDiskType, r.WorkerDiskConfig.BootDiskSize, poolSize))

	if r.WorkerDiskConfig.NumLocalSSDs > 0 {
		costComponents = append(costComponents, scratchDiskCostComponent(r.Region, r.WorkerConfig.PurchaseOption, int(r.WorkerDiskConfig.NumLocalSSDs)))
	}

	//Preemptible Worker Node Config
	costComponents = append(costComponents, computeCostDataProcComponent(r.Region, r.PreemptibleWorkerConfig.MachineType, r.PreemptibleWorkerConfig.PurchaseOption, 1, nil,"Preemptible Worker Node"))

	costComponents = append(costComponents, computeDiskCostComponent(r.Region, r.PreemptibleWorkerDiskConfig.BootDiskType, r.PreemptibleWorkerDiskConfig.BootDiskSize, poolSize))

	if r.WorkerDiskConfig.NumLocalSSDs > 0 {
		costComponents = append(costComponents, scratchDiskCostComponent(r.Region, r.PreemptibleWorkerConfig.PurchaseOption, int(r.PreemptibleWorkerDiskConfig.NumLocalSSDs)))
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    DataprocClusterUsageSchema,
		CostComponents: costComponents,
	}
}

func (r *DataprocCluster) instanceCostComponent() *schema.CostComponent {
	return &schema.CostComponent{
		Name:           fmt.Sprintf("Instance (on-demand, %s)", r.InstanceType),
		Unit:           "hours",
		UnitMultiplier: decimal.NewFromInt(1),
		// InstanceCount goes into the hourly quantity to increase the price with the number
		// of instances that are provisioned.
		HourlyQuantity: decimalPtr(decimal.NewFromInt(r.InstanceCount)),
		// ProductFilters find the actual price from the Infracost pricing database.
		// see https://github.com/infracost/infracost/blob/master/CONTRIBUTING.md#finding-prices
		// for more info on how to generate these filters to fetch the prices you need.
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("google"),
			Region:        strPtr(r.Region),
			Service:       strPtr("My google Service"),
			ProductFamily: strPtr("My google Resource family"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "usagetype", ValueRegex: regexPtr(fmt.Sprintf("^%s$", r.UsageType))},
				{Key: "instanceType", Value: strPtr(r.InstanceType)},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption: strPtr("on_demand"),
		},
	}
}

func (r *DataprocCluster) dataProcessedCostComponent() *schema.CostComponent {
	return &schema.CostComponent{
		Name:           "Data processed",
		Unit:           "GB",
		UnitMultiplier: decimal.NewFromInt(1),
		HourlyQuantity: floatPtrToDecimalPtr(r.MonthlyDataProcessedGB),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("google"),
			Region:        strPtr(r.Region),
			Service:       strPtr("My google Service"),
			ProductFamily: strPtr("My google Resource family"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "usagetype", Value: strPtr("UsageBytes")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption: strPtr("on_demand"),
		},
	}
}
