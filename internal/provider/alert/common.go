package alert

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

const (
	description = "Alerts let you monitor data including network traffic and activities, or various security events like password resets and missing certificates.\n" +
		"You can examine and filter any type of event, as well as define alert notifications to be sent to email, webhooks (integrating with SaaS apps), PagerDuty or Slack.\n" +
		"Alerts can be configured using either `spike_condition` or `threshold_condition`"
	channelsDesc      = "List of notification channel IDs."
	groupByDesc       = "The group by field name."
	notifyMessageDesc = "Creates a custom message that will be sent to your notification channels.\n" +
		"	You can use free text and/or alert field names surrounded with a \"${ }\". For example, \"${hits} have failed to login\"."
	sourceTypeDesc = "Logs type. Supported log types:\n" +
		"	- **security_audit**- The `security_audit` logs provide the administrator visibility into events which are generated by device and user security-related activity, such as user authenticating into Proofpoint NaaS, users changing their passwords, posture check failures, etc. See [here](https://help.metanetworks.com/knowledgebase/admin_console_logs/#security-logs) for details.\n" +
		"	- **api_audit** - The `api_audit` logs capture details of administrator activity: the timestamp and identity of administrators who accessed the Proofpoint NaaS tenant, and configuration changes that were made by the administrator. See [here](https://help.metanetworks.com/knowledgebase/admin_console_logs/#audit-logs) for details.\n" +
		"	- **traffic_audit** - The `traffic_audit` logs provide detailed visibility into each element in the system covering network traffic including DNS and other OSI Layer 3 and 4 traffic details. See [here](https://help.metanetworks.com/knowledgebase/admin_console_logs/#traffic-logs) for details.\n" +
		"	- **metaproxy_audit** - The `metaproxy_audit` logs provide the administrator visibility into the clientless access of their employees to web applications configured via EasyLink policy. See [here](https://help.metanetworks.com/knowledgebase/admin_console_logs/#metaconnect-web-logs) for details.\n" +
		"	- **webfilter_audit** - The `webfilter_audit` logs provide the administrator visibility into the events generated by the Web Security engine. See [here](https://help.metanetworks.com/knowledgebase/logs_ws/) for details\n."
	minHitsDesc    = "Minimum number of hits in current window to check the spike."
	spikeRatioDesc = "The difference between hits that triggers alert (in percents)."
	spikeTypeDesc  = "Spike type, ENUM: `up`, `down`, `both`."
	timeDiffDesc   = "Time difference in minutes between current and reference window, Enum: `1`, `3`, `5`, `60`, `1440`, `10080`."
	formulaDesc    = "Mathematical formula to run on the events, ENUM: `count`."
	opDesc         = "Operator used to compare to the threshold, ENUM: `greater`, `greaterequals`, `less`, `lessequals`, `equals`."
	thresholdDesc  = "The threshold to compare result of the formula."
	windowDesc     = "The time window of the check (in mins), ENUM: `1`, `3`, `5`, `10`, `30`, `60`, `360`, `1440`, `2880`, `10080`."
)

var excludedKeys = []string{"id", "spike_condition", "threshold_condition"}

func alertToResource(d *schema.ResourceData, a *client.Alert) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId(a.ID)
	err := client.MapResponseToResource(a, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	if a.SpikeCondition != nil {
		SpikeConditionToResource := []map[string]interface{}{
			{
				"min_hits":    a.SpikeCondition.MinHits,
				"spike_ratio": a.SpikeCondition.SpikeRatio,
				"spike_type":  a.SpikeCondition.SpikeType,
				"time_diff":   a.SpikeCondition.TimeDiff,
			},
		}
		err = d.Set("spike_condition", SpikeConditionToResource)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if a.ThresholdCondition != nil {
		ThresholdConditionToResource := []map[string]interface{}{
			{
				"formula":   a.ThresholdCondition.Formula,
				"op":        a.ThresholdCondition.Op,
				"threshold": a.ThresholdCondition.Threshold,
			},
		}
		err = d.Set("threshold_condition", ThresholdConditionToResource)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return diags
}

func alertRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	a, err := client.GetAlert(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
			return
		} else {
			return diag.FromErr(err)
		}
	}
	return alertToResource(d, a)
}
func alertCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	body := client.NewAlert(d)
	a, err := client.CreateAlert(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return alertToResource(d, a)
}

func alertUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewAlert(d)
	a, err := client.UpdateAlert(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return alertToResource(d, a)
}

func alertDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteAlert(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return
}
