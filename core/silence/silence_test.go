package silence_test

import (
	"testing"

	"github.com/goto/siren/core/receiver"
	"github.com/goto/siren/core/silence"
	"github.com/goto/siren/core/subscription"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"
)

func TestSilence_Evaluate(t *testing.T) {
	failedCases := []struct {
		name      string
		silence   silence.Silence
		rcv       subscription.Receiver
		want      bool
		errString string
	}{
		{
			name: "silence type that is not subscription type would return error",
			silence: silence.Silence{
				ID:   "silence-id",
				Type: "test",
			},
			errString: "silence id 'silence-id' type is not subscription, type is 'test' instead",
		},
		{
			name: "rule that is not evaluated to boolean would return error",
			silence: silence.Silence{
				ID:       "silence-id",
				Type:     silence.TypeSubscription,
				TargetID: 12,
				TargetExpression: map[string]interface{}{
					"rule": "1 + 1",
				},
			},
			rcv: subscription.Receiver{
				ID:   12,
				Type: receiver.TypePagerDuty,
			},
			errString: "rule evaluation result is not boolean: 2",
		},
		{
			name: "rule that cannot be evaluated would return error",
			silence: silence.Silence{
				ID:       "silence-id",
				Type:     silence.TypeSubscription,
				TargetID: 12,
				TargetExpression: map[string]interface{}{
					"rule": "test",
				},
			},
			rcv: subscription.Receiver{
				ID:   12,
				Type: receiver.TypePagerDuty,
			},
			errString: "rule evaluation result is not boolean: <nil>",
		},
	}

	sucessCases := []struct {
		name      string
		silence   silence.Silence
		rcv       subscription.Receiver
		want      bool
		errString string
	}{
		{
			name: "match by empty rule would pass",
			silence: silence.Silence{
				Type:     silence.TypeSubscription,
				TargetID: 12,
				TargetExpression: map[string]interface{}{
					"rule": "",
				},
			},
			rcv: subscription.Receiver{
				ID:   12,
				Type: receiver.TypePagerDuty,
			},
			want: true,
		},
		{
			name: "no rule key in target expression would return empty string",
			silence: silence.Silence{
				ID:               "silence-id",
				Type:             silence.TypeSubscription,
				TargetID:         12,
				TargetExpression: map[string]interface{}{},
			},
			rcv: subscription.Receiver{
				ID:   12,
				Type: receiver.TypePagerDuty,
			},
			want: true,
		},
		{
			name: "match by `true` rule would pass",
			silence: silence.Silence{
				Type:     silence.TypeSubscription,
				TargetID: 12,
				TargetExpression: map[string]interface{}{
					"rule": "true",
				},
			},
			rcv: subscription.Receiver{
				ID:   12,
				Type: receiver.TypePagerDuty,
			},
			want: true,
		},
		{
			name: "match by receiver id and type would pass",
			silence: silence.Silence{
				Type:     silence.TypeSubscription,
				TargetID: 12,
				TargetExpression: map[string]interface{}{
					"rule": "(ID == 12) and (Type == 'pagerduty')",
				},
			},
			rcv: subscription.Receiver{
				ID:   12,
				Type: receiver.TypePagerDuty,
			},
			want: true,
		},
		{
			name: "match multiple receivers would pass",
			silence: silence.Silence{
				Type:     silence.TypeSubscription,
				TargetID: 12,
				TargetExpression: map[string]interface{}{
					"rule": "(ID == 12) or (ID == 16)",
				},
			},
			rcv: subscription.Receiver{
				ID:   12,
				Type: receiver.TypePagerDuty,
			},
			want: true,
		},
	}

	tests := append(sucessCases, failedCases...)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapSubscriptionReceiver := map[string]interface{}{}

			err := mapstructure.Decode(tt.rcv, &mapSubscriptionReceiver)
			require.NoError(t, err)

			got, err := tt.silence.EvaluateSubscriptionRule(mapSubscriptionReceiver)
			if err != nil {
				if err.Error() != tt.errString {
					t.Errorf("silence.Silence.Evaluate() error = %v, expected was %v", err, tt.errString)
				}
			}
			if got != tt.want {
				t.Errorf("silence.Silence.Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSilence_Validate(t *testing.T) {
	tests := []struct {
		name    string
		sil     silence.Silence
		wantErr bool
	}{
		{
			name: "should return error if type subscription and target id is empty or zero",
			sil: silence.Silence{
				Type: silence.TypeSubscription,
			},
			wantErr: true,
		},
		{
			name: "should return error if type labels and target expression is empty",
			sil: silence.Silence{
				Type: silence.TypeMatchers,
			},
			wantErr: true,
		},
		{
			name: "should return no error if type subscription and target id is not empty or zero",
			sil: silence.Silence{
				Type:     silence.TypeSubscription,
				TargetID: 1,
			},
		},
		{
			name: "should return error if type labels and target expression is not empty",
			sil: silence.Silence{
				Type: silence.TypeMatchers,
				TargetExpression: map[string]interface{}{
					"k1": "v1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.sil.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Silence.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
