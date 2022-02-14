package route53domains

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53domains"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceRegisteredDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceRegisteredDomainCreate,
		Read:   resourceRegisteredDomainRead,
		Update: resourceRegisteredDomainUpdate,
		Delete: resourceRegisteredDomainDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"abuse_contact_email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"abuse_contact_phone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"admin_privacy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"expiration_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name_server": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 6,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"glue_ips": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 2,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.IsIPAddress,
							},
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(1, 255),
								validation.StringMatch(regexp.MustCompile(`[a-zA-Z0-9_\-.]*`), "can contain only alphabetical characters (A-Z or a-z), numeric characters (0-9), underscore (_), the minus sign (-), and the period (.)"),
							),
						},
					},
				},
			},
			"registrant_privacy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"registrar_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"registrar_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reseller": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
			"tech_privacy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"updated_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"whois_server": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceRegisteredDomainCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53DomainsConn

	domainName := d.Get("domain_name").(string)
	domainDetail, err := findDomainDetailByName(conn, domainName)

	if err != nil {
		return fmt.Errorf("error reading Route 53 Domains Domain (%s): %w", domainName, err)
	}

	d.SetId(aws.StringValue(domainDetail.DomainName))

	if adminPrivacy, registrantPrivacy, techPrivacy := d.Get("admin_privacy").(bool), d.Get("registrant_privacy").(bool), d.Get("tech_privacy").(bool); adminPrivacy != aws.BoolValue(domainDetail.AdminPrivacy) || registrantPrivacy != aws.BoolValue(domainDetail.RegistrantPrivacy) || techPrivacy != aws.BoolValue(domainDetail.TechPrivacy) {
		if err := modifyDomainContactPrivacy(conn, d.Id(), adminPrivacy, registrantPrivacy, techPrivacy, d.Timeout(schema.TimeoutCreate)); err != nil {
			return err
		}
	}

	if v := d.Get("auto_renew").(bool); v != aws.BoolValue(domainDetail.AutoRenew) {
		if err := modifyDomainAutoRenew(conn, d.Id(), v); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("name_server"); ok && len(v.([]interface{})) > 0 {
		nameservers := expandNameservers(v.([]interface{}))

		if !reflect.DeepEqual(nameservers, domainDetail.Nameservers) {
			if err := modifyDomainNameservers(conn, d.Id(), nameservers, d.Timeout(schema.TimeoutCreate)); err != nil {
				return err
			}
		}
	}

	tags, err := ListTags(conn, d.Id())

	if err != nil {
		return fmt.Errorf("error listing tags for Route 53 Domains Domain (%s): %w", d.Id(), err)
	}

	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig
	newTags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{}))).IgnoreConfig(ignoreTagsConfig)
	oldTags := tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	if !oldTags.Equal(newTags) {
		if err := UpdateTags(conn, d.Id(), oldTags, newTags); err != nil {
			return fmt.Errorf("error updating Route 53 Domains Domain (%s) tags: %w", d.Id(), err)
		}
	}

	return resourceRegisteredDomainRead(d, meta)
}

func resourceRegisteredDomainRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53DomainsConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	domainDetail, err := findDomainDetailByName(conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Route 53 Domains Domain %s not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading Route 53 Domains Domain (%s): %w", d.Id(), err)
	}

	d.Set("abuse_contact_email", domainDetail.AbuseContactEmail)
	d.Set("abuse_contact_phone", domainDetail.AbuseContactPhone)
	d.Set("admin_privacy", domainDetail.AdminPrivacy)
	d.Set("auto_renew", domainDetail.AutoRenew)
	if domainDetail.CreationDate != nil {
		d.Set("creation_date", aws.TimeValue(domainDetail.CreationDate).Format(time.RFC3339))
	} else {
		d.Set("creation_date", nil)
	}
	d.Set("domain_name", domainDetail.DomainName)
	if domainDetail.ExpirationDate != nil {
		d.Set("expiration_date", aws.TimeValue(domainDetail.ExpirationDate).Format(time.RFC3339))
	} else {
		d.Set("expiration_date", nil)
	}
	if err := d.Set("name_server", flattenNameservers(domainDetail.Nameservers)); err != nil {
		return fmt.Errorf("error setting name_servers: %w", err)
	}
	d.Set("registrant_privacy", domainDetail.RegistrantPrivacy)
	d.Set("registrar_name", domainDetail.RegistrarName)
	d.Set("registrar_url", domainDetail.RegistrarUrl)
	d.Set("reseller", domainDetail.Reseller)
	d.Set("status_list", aws.StringValueSlice(domainDetail.StatusList))
	d.Set("tech_privacy", domainDetail.TechPrivacy)
	if domainDetail.UpdatedDate != nil {
		d.Set("updated_date", aws.TimeValue(domainDetail.UpdatedDate).Format(time.RFC3339))
	} else {
		d.Set("updated_date", nil)
	}
	d.Set("whois_server", domainDetail.WhoIsServer)

	tags, err := ListTags(conn, d.Id())

	if err != nil {
		return fmt.Errorf("error listing tags for Route 53 Domains Domain (%s): %w", d.Id(), err)
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	return nil
}

func resourceRegisteredDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53DomainsConn

	if d.HasChanges("admin_privacy", "registrant_privacy", "tech_privacy") {
		if err := modifyDomainContactPrivacy(conn, d.Id(), d.Get("admin_privacy").(bool), d.Get("registrant_privacy").(bool), d.Get("tech_privacy").(bool), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}

	if d.HasChange("auto_renew") {
		if err := modifyDomainAutoRenew(conn, d.Id(), d.Get("auto_renew").(bool)); err != nil {
			return err
		}
	}

	if d.HasChange("name_server") {
		if v, ok := d.GetOk("name_server"); ok && len(v.([]interface{})) > 0 {
			if err := modifyDomainNameservers(conn, d.Id(), expandNameservers(v.([]interface{})), d.Timeout(schema.TimeoutUpdate)); err != nil {
				return err
			}
		}
	}

	// TODO TransferLock
	// TODO Contacts

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Id(), o, n); err != nil {
			return fmt.Errorf("error updating Route 53 Domains Domain (%s) tags: %w", d.Id(), err)
		}
	}

	return resourceRegisteredDomainRead(d, meta)
}

func resourceRegisteredDomainDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Route 53 Domains Registered Domain (%s) not deleted, removing from state", d.Id())

	return nil
}

func modifyDomainAutoRenew(conn *route53domains.Route53Domains, domainName string, autoRenew bool) error {
	if autoRenew {
		input := &route53domains.EnableDomainAutoRenewInput{
			DomainName: aws.String(domainName),
		}

		log.Printf("[DEBUG] Enabling Route 53 Domains Domain auto-renew: %s", input)
		_, err := conn.EnableDomainAutoRenew(input)

		if err != nil {
			return fmt.Errorf("error enabling Route 53 Domains Domain (%s) auto-renew: %w", domainName, err)
		}
	} else {
		input := &route53domains.DisableDomainAutoRenewInput{
			DomainName: aws.String(domainName),
		}

		log.Printf("[DEBUG] Disabling Route 53 Domains Domain auto-renew: %s", input)
		_, err := conn.DisableDomainAutoRenew(input)

		if err != nil {
			return fmt.Errorf("error disabling Route 53 Domains Domain (%s) auto-renew: %w", domainName, err)
		}
	}

	return nil
}

func modifyDomainContactPrivacy(conn *route53domains.Route53Domains, domainName string, adminPrivacy, registrantPrivacy, techPrivacy bool, timeout time.Duration) error {
	input := &route53domains.UpdateDomainContactPrivacyInput{
		AdminPrivacy:      aws.Bool(adminPrivacy),
		DomainName:        aws.String(domainName),
		RegistrantPrivacy: aws.Bool(registrantPrivacy),
		TechPrivacy:       aws.Bool(techPrivacy),
	}

	log.Printf("[DEBUG] Updating Route 53 Domains Domain contact privacy: %s", input)
	output, err := conn.UpdateDomainContactPrivacy(input)

	if err != nil {
		return fmt.Errorf("error enabling Route 53 Domains Domain (%s) contact privacy: %w", domainName, err)
	}

	if _, err := waitOperationSucceeded(conn, aws.StringValue(output.OperationId), timeout); err != nil {
		return fmt.Errorf("error waiting for Route 53 Domains Domain (%s) contact privacy update: %w", domainName, err)
	}

	return nil
}

func modifyDomainNameservers(conn *route53domains.Route53Domains, domainName string, nameservers []*route53domains.Nameserver, timeout time.Duration) error {
	input := &route53domains.UpdateDomainNameserversInput{
		DomainName:  aws.String(domainName),
		Nameservers: nameservers,
	}

	log.Printf("[DEBUG] Updating Route 53 Domains Domain name servers: %s", input)
	output, err := conn.UpdateDomainNameservers(input)

	if err != nil {
		return fmt.Errorf("error updating Route 53 Domains Domain (%s) name servers: %w", domainName, err)
	}

	if _, err := waitOperationSucceeded(conn, aws.StringValue(output.OperationId), timeout); err != nil {
		return fmt.Errorf("error waiting for Route 53 Domains Domain (%s) name servers update: %w", domainName, err)
	}

	return nil
}

func findDomainDetailByName(conn *route53domains.Route53Domains, name string) (*route53domains.GetDomainDetailOutput, error) {
	input := &route53domains.GetDomainDetailInput{
		DomainName: aws.String(name),
	}

	output, err := conn.GetDomainDetail(input)

	if tfawserr.ErrMessageContains(err, route53domains.ErrCodeInvalidInput, "not found") {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

func findOperationDetailByID(conn *route53domains.Route53Domains, id string) (*route53domains.GetOperationDetailOutput, error) {
	input := &route53domains.GetOperationDetailInput{
		OperationId: aws.String(id),
	}

	output, err := conn.GetOperationDetail(input)

	if tfawserr.ErrMessageContains(err, route53domains.ErrCodeInvalidInput, "No operation found") {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

func statusOperation(conn *route53domains.Route53Domains, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := findOperationDetailByID(conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}

func waitOperationSucceeded(conn *route53domains.Route53Domains, id string, timeout time.Duration) (*route53domains.GetOperationDetailOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{route53domains.OperationStatusSubmitted, route53domains.OperationStatusInProgress},
		Target:  []string{route53domains.OperationStatusSuccessful},
		Timeout: timeout,
		Refresh: statusOperation(conn, id),
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*route53domains.GetOperationDetailOutput); ok {
		tfresource.SetLastError(err, errors.New(aws.StringValue(output.Message)))

		return output, err
	}

	return nil, err
}

func flattenNameserver(apiObject *route53domains.Nameserver) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.GlueIps; v != nil {
		tfMap["glue_ips"] = aws.StringValueSlice(v)
	}

	if v := apiObject.Name; v != nil {
		tfMap["name"] = aws.StringValue(v)
	}

	return tfMap
}

func expandNameserver(tfMap map[string]interface{}) *route53domains.Nameserver {
	if tfMap == nil {
		return nil
	}

	apiObject := &route53domains.Nameserver{}

	if v, ok := tfMap["glue_ips"].(*schema.Set); ok && v.Len() > 0 {
		apiObject.GlueIps = flex.ExpandStringSet(v)
	}

	if v, ok := tfMap["name"].(string); ok && v != "" {
		apiObject.Name = aws.String(v)
	}

	return apiObject
}

func expandNameservers(tfList []interface{}) []*route53domains.Nameserver {
	if len(tfList) == 0 {
		return nil
	}

	var apiObjects []*route53domains.Nameserver

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		apiObject := expandNameserver(tfMap)

		if apiObject == nil {
			continue
		}

		apiObjects = append(apiObjects, apiObject)
	}

	return apiObjects
}

func flattenNameservers(apiObjects []*route53domains.Nameserver) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		tfList = append(tfList, flattenNameserver(apiObject))
	}

	return tfList
}
