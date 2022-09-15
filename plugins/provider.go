package plugins

// For rule, provider plugin needs to fullfill this interface.
//
// type RuleUploader interface {
// 	UpsertRule(ctx context.Context, rl *Rule, templateToUpdate *template.Template, namespaceURN string) error
// }
//
// For subscription, provider plugin needs to fullfill this interface.
// Notice there are subscription and subscriptions struct passed. Siren passes the list of subscriptions of
// a specific namespace in case the provier plugin needed.T
//
// type ProviderPlugin interface {
// 	CreateSubscription(ctx context.Context, sub *Subscription, subs []Subscription, namespaceURN string) error
// 	UpdateSubscription(ctx context.Context, sub *Subscription, subs []Subscription, namespaceURN string) error
// 	DeleteSubscription(ctx context.Context, sub *Subscription, subs []Subscription, namespaceURN string) error
// }
