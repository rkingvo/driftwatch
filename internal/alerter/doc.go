// Package alerter dispatches notifications when configuration drift is
// detected between running containers and their source manifests.
//
// Usage:
//
//	log := alerter.NewLogNotifier(os.Stderr)
//	a   := alerter.New(log)
//	if err := a.Send(report); err != nil {
//	    log.Println("alert error:", err)
//	}
//
// Additional Notifier implementations (webhook, Slack, PagerDuty, etc.) can
// be registered by passing them to alerter.New.
package alerter
