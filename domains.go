package dnspod

import (
	"encoding/json"
	"fmt"
	"sort"
)

const (
	methodDomainList   = "Domain.List"
	methodDomainCreate = "Domain.Create"
	methodDomainInfo   = "Domain.Info"
	methodDomainRemove = "Domain.Remove"
	methodRecordLine   = "Record.Line"
)

// DomainInfo handles domain information.
type DomainInfo struct {
	DomainTotal   json.Number `json:"domain_total,omitempty"`
	AllTotal      json.Number `json:"all_total,omitempty"`
	MineTotal     json.Number `json:"mine_total,omitempty"`
	ShareTotal    json.Number `json:"share_total,omitempty"`
	VipTotal      json.Number `json:"vip_total,omitempty"`
	IsMarkTotal   json.Number `json:"ismark_total,omitempty"`
	PauseTotal    json.Number `json:"pause_total,omitempty"`
	ErrorTotal    json.Number `json:"error_total,omitempty"`
	LockTotal     json.Number `json:"lock_total,omitempty"`
	SpamTotal     json.Number `json:"spam_total,omitempty"`
	VipExpire     json.Number `json:"vip_expire,omitempty"`
	ShareOutTotal json.Number `json:"share_out_total,omitempty"`
}

// Domain handles domain.
type Domain struct {
	ID               json.Number `json:"id,omitempty"`
	Name             string      `json:"name,omitempty"`
	PunyCode         string      `json:"punycode,omitempty"`
	Grade            string      `json:"grade,omitempty"`
	GradeTitle       string      `json:"grade_title,omitempty"`
	Status           string      `json:"status,omitempty"`
	ExtStatus        string      `json:"ext_status,omitempty"`
	Records          string      `json:"records,omitempty"`
	GroupID          json.Number `json:"group_id,omitempty"`
	IsMark           string      `json:"is_mark,omitempty"`
	Remark           string      `json:"remark,omitempty"`
	IsVIP            string      `json:"is_vip,omitempty"`
	SearchenginePush string      `json:"searchengine_push,omitempty"`
	UserID           string      `json:"user_id,omitempty"`
	CreatedOn        string      `json:"created_on,omitempty"`
	UpdatedOn        string      `json:"updated_on,omitempty"`
	TTL              json.Number `json:"ttl,omitempty"`
	CNameSpeedUp     string      `json:"cname_speedup,omitempty"`
	Owner            string      `json:"owner,omitempty"`
	AuthToAnquanBao  bool        `json:"auth_to_anquanbao,omitempty"`
	NameServer       []string    `json:"dnspod_ns,omitempty"`
}

type DomainCreateResp struct {
	Id       string   `json:"id"`
	Punycode string   `json:"punycode"`
	Domain   string   `json:"domain"`
	GradeNs  []string `json:"grade_ns"`
}

type domainListWrapper struct {
	Status  Status     `json:"status"`
	Info    DomainInfo `json:"info"`
	Domains []Domain   `json:"domains"`
}

type domainWrapper struct {
	Status Status     `json:"status"`
	Info   DomainInfo `json:"info"`
	Domain Domain     `json:"domain"`
}

type domainCreateWrapper struct {
	Status Status           `json:"status"`
	Domain DomainCreateResp `json:"domain"`
}

// DomainsService handles communication with the domain related methods of the DNSPod API.
//
// DNSPod API docs:
// - https://www.dnspod.cn/docs/domains.html
// - https://docs.dnspod.com/api/
type DomainsService struct {
	client *Client
}

// List the domains.
//
// DNSPod API docs:
// - https://www.dnspod.cn/docs/domains.html#domain-list
// - https://docs.dnspod.com/api/5fe1b40a6e336701a2111f5b/
func (s *DomainsService) List() ([]Domain, *Response, error) {
	payload := s.client.CommonParams.toPayLoad()
	payload.Set("length", fmt.Sprintf("%d", 3000))

	var all []Domain
	times := 0
GETALLDOMAINS:
	payload.Set("offset", fmt.Sprintf("%d", times*3000))
	returnedDomains := domainListWrapper{}
	res, err := s.client.post(methodDomainList, payload, &returnedDomains)
	if err != nil {
		return nil, res, err
	}

	if returnedDomains.Status.Code != "1" {
		return nil, nil, fmt.Errorf("could not get domains: %s", returnedDomains.Status.Message)
	}
	all = append(all, returnedDomains.Domains...)
	total, err := returnedDomains.Info.AllTotal.Int64()
	if err != nil {
		return nil, nil, err
	}
	if int64(len(all)) < total {
		times++
		goto GETALLDOMAINS
	}

	return returnedDomains.Domains, res, nil
}

// Create a new domain.
//
// DNSPod API docs:
// - https://www.dnspod.cn/docs/domains.html#domain-create
// - https://docs.dnspod.com/api/5fe1a9e36e336701a2111d3d/
func (s *DomainsService) Create(domainAttributes Domain) (DomainCreateResp, *Response, error) {
	payload := s.client.CommonParams.toPayLoad()
	payload.Set("domain", domainAttributes.Name)
	payload.Set("group_id", domainAttributes.GroupID.String())
	payload.Set("is_mark", domainAttributes.IsMark)

	returnedDomain := domainCreateWrapper{}

	res, err := s.client.post(methodDomainCreate, payload, &returnedDomain)
	if err != nil {
		return DomainCreateResp{}, res, err
	}

	if returnedDomain.Status.Code != "1" {
		return DomainCreateResp{}, nil, fmt.Errorf("code: %s,message: %s", returnedDomain.Status.Code, returnedDomain.Status.Message)
	}

	return returnedDomain.Domain, res, nil
}

// Get fetches a domain.
//
// DNSPod API docs:
// - https://www.dnspod.cn/docs/domains.html#domain-info
// - https://docs.dnspod.com/api/5fe1b37d6e336701a2111f2b/
func (s *DomainsService) Get(domainId string) (Domain, *Response, error) {
	payload := s.client.CommonParams.toPayLoad()
	payload.Set("domain_id", domainId)

	returnedDomain := domainWrapper{}

	res, err := s.client.post(methodDomainInfo, payload, &returnedDomain)
	if err != nil {
		return Domain{}, res, err
	}

	return returnedDomain.Domain, res, nil
}

// Delete a domain.
//
// DNSPod API docs:
// - https://dnsapi.cn/Domain.Remove
// - https://docs.dnspod.com/api/5fe1ac446e336701a2111dd1/
func (s *DomainsService) Delete(domainId string) (*Response, error) {
	payload := s.client.CommonParams.toPayLoad()
	payload.Set("domain_id", domainId)

	returnedDomain := domainWrapper{}

	res, err := s.client.post(methodDomainRemove, payload, &returnedDomain)
	if err != nil {
		return nil, err
	}

	if returnedDomain.Status.Code != "1" {
		return nil, fmt.Errorf("code: %s, message: %s", returnedDomain.Status.Code, returnedDomain.Status.Message)
	}

	return res, nil
}

type Line struct {
	LineName string `json:"line_name"`
	LineId   string `json:"line_id"`
}

type LineInfo struct {
	Status  Status         `json:"status"`
	LineIds map[string]any `json:"line_ids"`
}

// GetLines
//
// get lines of record which group by grade of domain
// valid grade: D_Free,D_Plus,D_Extra,D_Expert,D_Ultra,DP_Free,DP_Plus,DP_Extra,DP_Expert,DP_Ultra
func (s *DomainsService) GetLines(domain, domainGrade string) ([]Line, *Response, error) {
	grade := []string{
		"D_Free", "D_Plus", "D_Extra", "D_Expert", "D_Ultra",
		"DP_Free", "DP_Plus", "DP_Extra", "DP_Expert", "DP_Ultra",
		"DPG_Free", "DPG_Plus", "DPG_Extra", "DPG_Expert", "DPG_Ultra",
	}
	ok := false
	for _, i := range grade {
		if i == domainGrade {
			ok = true
		}
	}
	if !ok {
		return nil, nil, fmt.Errorf("invalid grade of domain: %s", domainGrade)
	}

	payload := s.client.CommonParams.toPayLoad()
	if s.client.CommonParams.IsInternational {
		if domain == "" {
			return nil, nil, fmt.Errorf("domain must need when connect to internatinal version")
		}
		payload.Set("domain", domain)
	}
	payload.Set("domain_grade", domainGrade)

	returnedLines := LineInfo{}

	res, err := s.client.post(methodRecordLine, payload, &returnedLines)
	if err != nil {
		return nil, nil, err
	}

	if returnedLines.Status.Code != "1" {
		return nil, nil, fmt.Errorf("code: %s, message: %s", returnedLines.Status.Code, returnedLines.Status.Message)
	}

	var items []Line
	for k, v := range returnedLines.LineIds {
		item := Line{LineName: k}

		switch v.(type) {
		case float64:
			item.LineId = fmt.Sprintf("%d", int(v.(float64)))
		case string:
			item.LineId = v.(string)
		}

		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].LineId < items[j].LineId
	})

	return items, res, nil
}
