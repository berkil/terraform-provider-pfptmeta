---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pfptmeta_easylink Resource - terraform-provider-pfptmeta"
subcategory: "Network Resources"
description: |-
  Proofpoint provides two access modes for end users:
  Agent-based, allows users access to resources natively, after establishing a VPN connection and authenticating using the Proofpoint Agent client.Clientless, or MetaConnect (MC), allows users to access resources via supported browsers, without any client installation.
  MetaConnect eliminates the need to install an agent and establish a VPN connection. Users access internal applications from a dedicated web (MetaConnect) portal. After authentication, the users see the list of applications that they are allowed to access. Alternatively, users can use a static FQDN to access their private apps directly. MetaConnect can be used to access web applications (HTTP or HTTPS), or servers via RDP, SSH or VNC.
  There are several use cases for using MetaConnect access mode. It is well suited for instances when an agent cannot be used, such as with external contractors or personal/BYO devices.
  Clientless applications are defined using EasyLinks. An EasyLink defines the application (internal host, protocol and port), the users assigned to this application, the URL type, etc.
  See here https://help.metanetworks.com/knowledgebase/easylinks/ for more details.
---

# pfptmeta_easylink (Resource)

Proofpoint provides two access modes for end users:

- Agent-based, allows users access to resources natively, after establishing a VPN connection and authenticating using the Proofpoint Agent client.

- Clientless, or MetaConnect (MC), allows users to access resources via supported browsers, without any client installation.

MetaConnect eliminates the need to install an agent and establish a VPN connection. Users access internal applications from a dedicated web (MetaConnect) portal. After authentication, the users see the list of applications that they are allowed to access. Alternatively, users can use a static FQDN to access their private apps directly. MetaConnect can be used to access web applications (HTTP or HTTPS), or servers via RDP, SSH or VNC.

There are several use cases for using MetaConnect access mode. It is well suited for instances when an agent cannot be used, such as with external contractors or personal/BYO devices.

Clientless applications are defined using EasyLinks. An EasyLink defines the application (internal host, protocol and port), the users assigned to this application, the URL type, etc.

See [here](https://help.metanetworks.com/knowledgebase/easylinks/) for more details.

## Example Usage

```terraform
resource "pfptmeta_group" "new_group" {
  name = "easylink-group"
}

locals {
  hostname = "test.example.com"
  ipv4     = "196.10.10.1"
}

resource "pfptmeta_certificate" "cert" {
  name = "certificate name"
  sans = [local.hostname]
}

resource "pfptmeta_network_element" "mapped-service" {
  name           = "mapped service name"
  mapped_service = local.ipv4
}

resource "pfptmeta_network_element_alias" "alias" {
  network_element_id = pfptmeta_network_element.mapped-service.id
  alias              = local.hostname
}

resource "pfptmeta_easylink" "meta_easylink" {
  name        = "meta easylink name"
  description = "meta easylink description"
  domain_name = local.hostname
  access_type = "meta"
  port        = 443
  protocol    = "https"
  viewers     = [pfptmeta_group.new_group.id]
}

resource "pfptmeta_easylink" "meta_rdp_easylink" {
  name        = "meta_rdp easylink name"
  description = "meta_rdp easylink description"
  domain_name = local.hostname
  access_type = "meta"
  port        = 3389
  protocol    = "rdp"
  viewers     = [pfptmeta_group.new_group.id]
  rdp {
    security               = "nla"
    server_keyboard_layout = "french"
  }
}

resource "pfptmeta_easylink" "redirect_easylink" {
  name              = "redirect easylink name"
  description       = "redirect easylink description"
  domain_name       = local.ipv4
  access_fqdn       = local.hostname
  access_type       = "redirect"
  port              = 443
  protocol          = "https"
  mapped_element_id = pfptmeta_network_element.mapped-service.id
  viewers           = [pfptmeta_group.new_group.id]
  certificate_id    = pfptmeta_certificate.cert.id
  root_path         = "/application"
}

resource "pfptmeta_easylink" "native_easylink" {
  name              = "native easylink name"
  description       = "native easylink description"
  domain_name       = local.ipv4
  access_fqdn       = local.hostname
  access_type       = "native"
  port              = 443
  protocol          = "https"
  mapped_element_id = pfptmeta_network_element.mapped-service.id
  viewers           = [pfptmeta_group.new_group.id]
  certificate_id    = pfptmeta_certificate.cert.id
  root_path         = "/application"
  proxy {
    rewrite_content_types = ["json", "html"]
    rewrite_http          = true
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **access_type** (String) When creating an Easylink, you need to select the appropriate access (URL) type.

	- **meta** – Use Proofpoint-generated URL as the application entry point and throughout the browsing session.

	- **redirect** – Use your own (vanity) FQDN as the application access entry point, and redirect to Proofpoint-generated URL for the rest of browsing session.

	- **native** – Use your own (vanity) FQDN as the application entry point and throughout the browsing session. Only web-based apps (HTTP or HTTPS) are supported.
- **domain_name** (String) FQDN or IPv4 of the application defined by the EasyLink.
- **name** (String)
- **port** (Number) The port of the application defined by the EasyLink.
- **protocol** (String) The protocol of the application defined by the EasyLink. ENUM: `ssh`, `rdp`, `vnc`, `http`, `https`.
- **viewers** (Set of String) User or group IDs that will be granted access to the application defined by the EasyLink.

### Optional

- **access_fqdn** (String) External FQDN to be associated with the current EasyLink, required when `access_type` is set to `redirect` or `native`.
- **audit** (Boolean) When enabled, all web traffic is logged to the MetaConnect Web log. Logging is only applicable when `protocol` is either `http` or `https`.
- **certificate_id** (String) Required when `access_type` is set to `redirect` or `native`. MetaConnect provides a secure connection with HTTPS. For the end-user browser to trust the domain, you must generate an SSL certificate for the external application FQDN using the `pfptmeta_certificate` resource.
- **description** (String)
- **enable_sni** (Boolean) Defines whether to enable SNI or not. The SNI can be enabled only when `protocol` is set to `https`.
- **mapped_element_id** (String) Hosting resource for Mapped Subnet or Mapped Service network elements if the host is to reside permanently within this resource. This field is required when the host is an IPv4 address.
- **proxy** (Block List, Max: 1) Additional proxy configuration, available only when `protocol` is set to `http` or `https`. (see [below for nested schema](#nestedblock--proxy))
- **rdp** (Block List, Max: 1) Additional RDP configuration, available only when `protocol` is set to `rdp`. (see [below for nested schema](#nestedblock--rdp))
- **root_path** (String) The root path of the application defined by the EasyLink, when `protocol` is `http` or `https`.

### Read-Only

- **id** (String) The ID of this resource.
- **version** (Number)

<a id="nestedblock--proxy"></a>
### Nested Schema for `proxy`

Optional:

- **enterprise_access** (Boolean) When enabled, it resets the session on source IP change to minimize latency if the new source IP has enterprise access to the EasyLink. Allowed only for default ports (80, 443) and when `access_type` is set to `redirect`.
- **hosts** (List of String) Additional hosts to be routed to the EasyLink.
- **http_host_header** (String) An overwrite to the HTTP host header. It is set to the value of `access_fqdn` when `access_type` is set to `native` and not allowed.
- **rewrite_content_types** (List of String) Response content types to be rewritten. ENUM: `html`, `json`, `javascript`, `text`. It is required when `rewrite_hosts` or `rewrite_http` are configured.
- **rewrite_hosts** (Boolean) Defines whether to rewrite hosts in the proxy response to the EasyLink host or not. Rewrites in responses with content type specified in `rewrite_content_types`.
- **rewrite_hosts_client** (Boolean) Selects whether to overwrite hosts in all browser client requests or not.
- **rewrite_http** (Boolean) Defines whether to overwrite all `http://` links in proxy response to `https://` or not. Rewrites in responses with content type specified in `rewrite_content_types`.
- **shared_cookies** (Boolean) Selects whether to share cookies between EasyLinks in the same region.


<a id="nestedblock--rdp"></a>
### Nested Schema for `rdp`

Optional:

- **remote_app** (String) The remote application to start on the remote desktop. If supported by your remote desktop server, only this application will be visible to the user.
- **remote_app_cmd_args** (String) The command-line arguments, if any, for the remote application.
- **remote_app_work_dir** (String) The working directory, if any, for the remote application.
- **security** (String) Dictates how data is encrypted and what type of authentication is performed. ENUM: `nla`, `rdp`.
- **server_keyboard_layout** (String) Server-supported keyboard layout. Enum: `english-us`, `german`, `french`, `swiss-french`, `italian`, `japanese`, `swedish`, `unicode`.
