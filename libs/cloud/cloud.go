package cloud

import (
	"flag"
	"reflect"
	"strings"
	"time"

	"github.com/onaio/sre-tooling/libs/cli/flags"
)

type Provider interface {
	getName() string
	getAllResources(filter *Filter) ([]*Resource, error)
}

type Resource struct {
	Provider     string
	ID           string
	Location     string
	ResourceType string
	LaunchTime   time.Time
	Tags         map[string]string
}

type Filter struct {
	Providers     []string
	ResourceTypes []string
	Regions       []string
	Tags          map[string]string
}

const tagFlagSeparator = ":"

func GetAllCloudResources(filter *Filter) ([]*Resource, error) {
	allResources := []*Resource{}

	aws := AWS{}
	if considerProvider(aws, filter) {
		awsResources, awsErr := aws.getAllResources(filter)
		if awsErr != nil {
			return nil, awsErr
		}
		allResources = append(allResources, awsResources...)
	}

	return allResources, nil
}

func GetTagKeys(resource *Resource) []string {
	keyObjects := reflect.ValueOf(resource.Tags).MapKeys()
	keys := make([]string, len(keyObjects))
	for i := 0; i < len(keyObjects); i++ {
		keys[i] = keyObjects[i].String()
	}

	return keys
}

func AddFlags(flagSet *flag.FlagSet) (*flags.StringArray, *flags.StringArray, *flags.StringArray, *flags.StringArray) {
	providerFlag := new(flags.StringArray)
	flagSet.Var(providerFlag, "provider", "Name of provider to filter using")
	regionFlag := new(flags.StringArray)
	flagSet.Var(regionFlag, "region", "Name of a provider region to filter using")
	typeFlag := new(flags.StringArray)
	flagSet.Var(typeFlag, "type", "Resource type to filter using")
	tagFlag := new(flags.StringArray)
	flagSet.Var(tagFlag, "tag", "Resource tag to filter using. Use the format \"tagKey"+tagFlagSeparator+"tagValue\"")

	return providerFlag, regionFlag, typeFlag, tagFlag
}

func GetFiltersFromCommandFlags(providerFlag *flags.StringArray, regionFlag *flags.StringArray, typeFlag *flags.StringArray, tagFlag *flags.StringArray) *Filter {
	filter := Filter{}
	if len(*providerFlag) > 0 {
		filter.Providers = *providerFlag
	}
	if len(*typeFlag) > 0 {
		filter.ResourceTypes = *typeFlag
	}
	if len(*regionFlag) > 0 {
		filter.Regions = *regionFlag
	}
	if len(*tagFlag) > 0 {
		for _, curTagPair := range *tagFlag {
			curKeyValue := strings.Split(curTagPair, tagFlagSeparator)
			if len(curKeyValue) == 2 {
				if filter.Tags == nil {
					filter.Tags = make(map[string]string)
				}
				filter.Tags[curKeyValue[0]] = curKeyValue[1]
			}
		}
	}

	return &filter
}

func considerProvider(providerIface interface{}, filter *Filter) bool {
	provider := providerIface.(Provider)
	if len(filter.Providers) == 0 {
		return true
	}

	for _, curProviderName := range filter.Providers {
		if strings.ToLower(curProviderName) == strings.ToLower(provider.getName()) {
			return true
		}
	}

	return false
}

func considerRegion(region string, filter *Filter) bool {
	if len(filter.Regions) == 0 {
		return true
	}

	for _, curRegion := range filter.Regions {
		if strings.ToLower(curRegion) == strings.ToLower(region) {
			return true
		}
	}

	return false
}

func considerResourceType(resourceType string, filter *Filter) bool {
	if len(filter.ResourceTypes) == 0 {
		return true
	}

	for _, curType := range filter.ResourceTypes {
		if strings.ToLower(curType) == strings.ToLower(resourceType) {
			return true
		}
	}

	return false
}

func considerTags(tags map[string]string, filter *Filter) bool {
	if len(filter.Tags) == 0 {
		return true
	}

	allOk := true
	for tagName, tagValue := range filter.Tags {
		if tagValueToCheck, ok := tags[tagName]; ok {
			if strings.ToLower(tagValue) != strings.ToLower(tagValueToCheck) {
				allOk = false
				break
			}
		} else {
			allOk = false
			break
		}
	}
	return allOk
}
