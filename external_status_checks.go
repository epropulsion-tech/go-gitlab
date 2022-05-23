package gitlab

import (
	"fmt"
	"net/http"
	"time"
)

// ExternalStatusChecksService handles communication with the external
// status check related methods of the GitLab API.
//
// GitLab API docs: https://docs.gitlab.com/ee/api/status_checks.html
type ExternalStatusChecksService struct {
	client *Client
}

type MergeStatusCheck struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ExternalURL string `json:"external_url"`
	Status      string `json:"status"`
}

type ProjectStatusCheck struct {
	ID                int                          `json:"id"`
	Name              string                       `json:"name"`
	ProjectID         int                          `json:"project_id"`
	ExternalURL       string                       `json:"external_url"`
	ProtectedBranches []StatusCheckProtectedBranch `json:"protected_branches"`
}

type StatusCheckProtectedBranch struct {
	ID                        int        `json:"id"`
	ProjectID                 int        `json:"project_id"`
	Name                      string     `json:"name"`
	CreatedAt                 *time.Time `json:"created_at"`
	UpdatedAt                 *time.Time `json:"updated_at"`
	CodeOwnerApprovalRequired bool       `json:"code_owner_approval_required"`
}

// ListMergeStatusChecks lists the external status checks that apply to it
// and their status for a single merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/ee/api/status_checks.html#list-status-checks-for-a-merge-request
func (s *ExternalStatusChecksService) ListMergeStatusChecks(pid interface{}, mr int, opt *ListOptions, options ...RequestOptionFunc) ([]*MergeStatusCheck, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/merge_requests/%d/status_checks", PathEscape(project), mr)

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var mscs []*MergeStatusCheck
	resp, err := s.client.Do(req, &mscs)
	if err != nil {
		return nil, resp, err
	}

	return mscs, resp, err
}

// ListProjectStatusChecks lists the project external status checks.
//
// GitLab API docs:
// https://docs.gitlab.com/ee/api/status_checks.html#get-project-external-status-checks
func (s *ExternalStatusChecksService) ListProjectStatusChecks(pid interface{}, opt *ListOptions, options ...RequestOptionFunc) ([]*ProjectStatusCheck, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/external_status_checks", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var pscs []*ProjectStatusCheck
	resp, err := s.client.Do(req, &pscs)
	if err != nil {
		return nil, resp, err
	}

	return pscs, resp, err
}

type SetExternalStatusCheckStatusOptions struct {
	SHA                   *string `url:"sha" json:"sha"`
	ExternalStatusCheckID *int    `url:"external_status_check_id" json:"external_status_check_id"`
	Status                *string `url:"status,omitempty" json:"status,omitempty"`
}

// SetExternalStatusCheckStatus set status of an external status check
//
// Gitlab API docs:
// https://docs.gitlab.com/ee/api/status_checks.html#set-status-of-an-external-status-check
func (s *ExternalStatusChecksService) SetExternalStatusCheckStatus(pid interface{}, mergeRequestIID int, opt *SetExternalStatusCheckStatusOptions, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/merge_requests/%d/status_check_responses", PathEscape(project), mergeRequestIID)

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

type CreateExternalStatusCheck struct {
	Name               *string `url:"name" json:"name"`
	ExternalURL        *string `url:"external_url" json:"external_url"`
	ProtectedBranchIDs *[]int  `url:"protected_branch_ids,omitempty" json:"protected_branch_ids,omitempty"`
}

func (s *ExternalStatusChecksService) CreateExternalStatusCheck(pid interface{}, opt *CreateExternalStatusCheck, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/external_status_checks", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

func (s *ExternalStatusChecksService) DeleteExternalStatusCheck(pid interface{}, checkID int, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}

	u := fmt.Sprintf("projects/%s/external_status_checks/%d", PathEscape(project), checkID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

type UpdateExternalStatusCheckOptions struct {
	Name               *string `url:"name,omitempty" json:"name,omitempty"`
	ExternalURL        *string `url:"external_url,omitempty" json:"external_url,omitempty"`
	ProtectedBranchIDs *[]int  `url:"protected_branch_ids,omitempty" json:"protected_branch_ids,omitempty"`
}

func (s *ExternalStatusChecksService) UpdateExternalStatusCheck(pid interface{}, checkID int, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}

	u := fmt.Sprintf("projects/%s/external_status_checks/%d", PathEscape(project), checkID)

	req, err := s.client.NewRequest(http.MethodPut, u, &UpdateExternalStatusCheckOptions{}, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
