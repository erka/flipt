package flipt

import "context"

type Requester interface {
	Request(ctx context.Context) Request
}

// Resource represents what resource or parent resource is being acted on.
type Resource string

// Subject returns the subject of the request.
type Subject string

// Action represents the action being taken on the resource.
type Action string

// Status represents the status of the request.
type Status string

const (
	ResourceNamespace      Resource = "namespace"
	ResourceFlag           Resource = "flag"
	ResourceSegment        Resource = "segment"
	ResourceAuthentication Resource = "authentication"
	ResourceOFREP          Resource = "ofrep"

	SubjectConstraint   Subject = "constraint"
	SubjectDistribution Subject = "distribution"
	SubjectFlag         Subject = "flag"
	SubjectNamespace    Subject = "namespace"
	SubjectRollout      Subject = "rollout"
	SubjectRule         Subject = "rule"
	SubjectSegment      Subject = "segment"
	SubjectToken        Subject = "token"
	SubjectVariant      Subject = "variant"

	ActionCreate   Action = "create"
	ActionDelete   Action = "delete"
	ActionUpdate   Action = "update"
	ActionRead     Action = "read"
	ActionEvaluate Action = "evaluate"

	StatusSuccess Status = "success"
	StatusDenied  Status = "denied"
)

type Request struct {
	Namespace string   `json:"namespace"`
	Resource  Resource `json:"resource"`
	Subject   Subject  `json:"subject"`
	Action    Action   `json:"action"`
	Status    Status   `json:"status"`
}

func WithNamespace(ns string) func(*Request) {
	return func(r *Request) {
		if ns != "" {
			r.Namespace = ns
		}
	}
}

func WithStatus(s Status) func(*Request) {
	return func(r *Request) {
		r.Status = s
	}
}

func WithSubject(s Subject) func(*Request) {
	return func(r *Request) {
		r.Subject = s
	}
}

func NewRequest(r Resource, a Action, opts ...func(*Request)) Request {
	req := Request{
		Resource:  r,
		Action:    a,
		Status:    StatusSuccess,
		Namespace: DefaultNamespace,
	}

	for _, opt := range opts {
		opt(&req)
	}

	return req
}

func newFlagScopedRequest(ns string, s Subject, a Action) Request {
	return NewRequest(ResourceFlag, a, WithNamespace(ns), WithSubject(s))
}

func newSegmentScopedRequest(ns string, s Subject, a Action) Request {
	return NewRequest(ResourceSegment, a, WithNamespace(ns), WithSubject(s))
}

// Namespaces
func (req *GetNamespaceRequest) Request(context.Context) Request {
	return NewRequest(ResourceNamespace, ActionRead, WithNamespace(req.Key))
}

func (req *ListNamespaceRequest) Request(context.Context) Request {
	return NewRequest(ResourceNamespace, ActionRead)
}

func (req *CreateNamespaceRequest) Request(context.Context) Request {
	return NewRequest(ResourceNamespace, ActionCreate, WithNamespace(req.Key))
}

func (req *UpdateNamespaceRequest) Request(context.Context) Request {
	return NewRequest(ResourceNamespace, ActionUpdate, WithNamespace(req.Key))
}

func (req *DeleteNamespaceRequest) Request(context.Context) Request {
	return NewRequest(ResourceNamespace, ActionDelete, WithNamespace(req.Key))
}

// Flags
func (req *GetFlagRequest) Request(context.Context) Request {
	return newFlagScopedRequest(req.NamespaceKey, SubjectFlag, ActionRead)
}

func (req *ListFlagRequest) Request(context.Context) Request {
	return newFlagScopedRequest(req.NamespaceKey, SubjectFlag, ActionRead)
}

func (req *CreateFlagRequest) Request(context.Context) Request {
	return newFlagScopedRequest(req.NamespaceKey, SubjectFlag, ActionCreate)
}

func (req *UpdateFlagRequest) Request(context.Context) Request {
	return newFlagScopedRequest(req.NamespaceKey, SubjectFlag, ActionUpdate)
}

func (req *DeleteFlagRequest) Request(context.Context) Request {
	return newFlagScopedRequest(req.NamespaceKey, SubjectFlag, ActionDelete)
}

// Variants
func (req *CreateVariantRequest) Request(context.Context) Request {
	return newFlagScopedRequest(req.NamespaceKey, SubjectVariant, ActionCreate)
}

func (req *UpdateVariantRequest) Request(context.Context) Request {
	return newFlagScopedRequest(req.NamespaceKey, SubjectVariant, ActionUpdate)
}

func (req *DeleteVariantRequest) Request(context.Context) Request {
	return newFlagScopedRequest(req.NamespaceKey, SubjectVariant, ActionDelete)
}

// Rules
func (req *ListRuleRequest) Request(context.Context) Request {
	return newFlagScopedRequest(req.NamespaceKey, SubjectRule, ActionRead)
}

func (req *GetRuleRequest) Request(context.Context) Request {
	return newFlagScopedRequest(req.NamespaceKey, SubjectRule, ActionRead)
}

func (req *CreateRuleRequest) Request(context.Context) Request {
	return newFlagScopedRequest(req.NamespaceKey, SubjectRule, ActionCreate)
}

func (req *UpdateRuleRequest) Request(context.Context) Request {
	return newFlagScopedRequest(req.NamespaceKey, SubjectRule, ActionUpdate)
}

func (req *OrderRulesRequest) Request(context.Context) Request {
	return newFlagScopedRequest(req.NamespaceKey, SubjectRule, ActionUpdate)
}

func (req *DeleteRuleRequest) Request(context.Context) Request {
	return newFlagScopedRequest(req.NamespaceKey, SubjectRule, ActionDelete)
}

// Rollouts
func (req *ListRolloutRequest) Request(context.Context) Request {
	return newFlagScopedRequest(req.NamespaceKey, SubjectRollout, ActionRead)
}

func (req *GetRolloutRequest) Request(context.Context) Request {
	return newFlagScopedRequest(req.NamespaceKey, SubjectRollout, ActionRead)
}

func (req *CreateRolloutRequest) Request(context.Context) Request {
	return newFlagScopedRequest(req.NamespaceKey, SubjectRollout, ActionCreate)
}

func (req *UpdateRolloutRequest) Request(context.Context) Request {
	return newFlagScopedRequest(req.NamespaceKey, SubjectRollout, ActionUpdate)
}

func (req *OrderRolloutsRequest) Request(context.Context) Request {
	return newFlagScopedRequest(req.NamespaceKey, SubjectRollout, ActionUpdate)
}

func (req *DeleteRolloutRequest) Request(context.Context) Request {
	return newFlagScopedRequest(req.NamespaceKey, SubjectRollout, ActionDelete)
}

// Segments
func (req *GetSegmentRequest) Request(context.Context) Request {
	return newSegmentScopedRequest(req.NamespaceKey, SubjectSegment, ActionRead)
}

func (req *ListSegmentRequest) Request(context.Context) Request {
	return newSegmentScopedRequest(req.NamespaceKey, SubjectSegment, ActionRead)
}

func (req *CreateSegmentRequest) Request(context.Context) Request {
	return newSegmentScopedRequest(req.NamespaceKey, SubjectSegment, ActionCreate)
}

func (req *UpdateSegmentRequest) Request(context.Context) Request {
	return newSegmentScopedRequest(req.NamespaceKey, SubjectSegment, ActionUpdate)
}

func (req *DeleteSegmentRequest) Request(context.Context) Request {
	return newSegmentScopedRequest(req.NamespaceKey, SubjectSegment, ActionDelete)
}

// Constraints
func (req *CreateConstraintRequest) Request(context.Context) Request {
	return newSegmentScopedRequest(req.NamespaceKey, SubjectConstraint, ActionCreate)
}

func (req *UpdateConstraintRequest) Request(context.Context) Request {
	return newSegmentScopedRequest(req.NamespaceKey, SubjectConstraint, ActionUpdate)
}

func (req *DeleteConstraintRequest) Request(context.Context) Request {
	return newSegmentScopedRequest(req.NamespaceKey, SubjectConstraint, ActionDelete)
}

// Distributions
func (req *CreateDistributionRequest) Request(context.Context) Request {
	return newSegmentScopedRequest(req.NamespaceKey, SubjectDistribution, ActionCreate)
}

func (req *UpdateDistributionRequest) Request(context.Context) Request {
	return newSegmentScopedRequest(req.NamespaceKey, SubjectDistribution, ActionUpdate)
}

func (req *DeleteDistributionRequest) Request(context.Context) Request {
	return newSegmentScopedRequest(req.NamespaceKey, SubjectDistribution, ActionDelete)
}
