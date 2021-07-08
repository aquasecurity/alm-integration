package router

import (
	"testing"
	"time"
)

var (
	singleRoute string = `
Name: tenant

routes:
- name: route1
  outputs: ["my-slack"]
  template: raw
  plugins:
   Policy-Show-All: true

templates:
- name: raw
  body: |
   package postee
   result:=input

outputs:
- name: my-slack
  type: slack
  enable: true
  url: https://hooks.slack.com/services/ABCDF/1234/TTT`
	noAssociatedOutput string = `
Name: tenant

routes:
- name: route1
  template: raw
  plugins:
   Policy-Show-All: true

templates:
- name: raw
  body: |
   package postee
   result:=input

outputs:
- name: my-slack
  type: slack
  enable: true
  url: https://hooks.slack.com/services/ABCDF/1234/TTT`
	twoRoutes string = `
Name: tenant

routes:
- name: route1
  outputs: ["my-slack"]
  template: raw
  plugins:
   Policy-Show-All: true

- name: route2
  outputs: ["my-slack"]
  template: raw
  plugins:
   Policy-Show-All: true

templates:
- name: raw
  body: |
   package postee
   result:=input

outputs:
- name: my-slack
  type: slack
  enable: true
  url: https://hooks.slack.com/services/ABCDF/1234/TTT`

	twoOutputs string = `
Name: tenant

routes:
- name: route1
  outputs: ["my-slack", "my-slack2"]
  template: raw
  plugins:
   Policy-Show-All: true

templates:
- name: raw
  body: |
   package postee
   result:=input

outputs:
- name: my-slack
  type: slack
  enable: true
  url: https://hooks.slack.com/services/ABCDF/1234/XXX
- name: my-slack2
  type: slack
  enable: true
  url: https://hooks.slack.com/services/ABCDF/1234/TTT`
	noOutputs string = `
Name: tenant

routes:
- name: route1
  outputs: ["my-slack3"]
  template: raw
  plugins:
   Policy-Show-All: true

templates:
- name: raw
  body: |
   package postee
   result:=input`
	noTemplates string = `
Name: tenant

routes:
- name: route1
  outputs: ["my-slack", "my-slack2"]
  template: raw
  plugins:
   Policy-Show-All: true

outputs:
- name: my-slack
  type: slack
  enable: true
  url: https://hooks.slack.com/services/ABCDF/1234/XXX
- name: my-slack2
  type: slack
  enable: true
  url: https://hooks.slack.com/services/ABCDF/1234/TTT`
	invalidTemplate string = `
Name: tenant

routes:
- name: route1
  outputs: ["my-slack"]
  template: rawx
  plugins:
   Policy-Show-All: true

templates:
- name: raw
  body: |
   package postee
   result:=input

outputs:
- name: my-slack
  type: slack
  enable: true
  url: https://hooks.slack.com/services/ABCDF/1234/TTT`
	invalidOutput string = `
Name: tenant

routes:
- name: route1
  outputs: ["x-slack"]
  template: raw
  plugins:
   Policy-Show-All: true

templates:
- name: raw
  body: |
   package postee
   result:=input

outputs:
- name: my-slack
  type: slack
  enable: true
  url: https://hooks.slack.com/services/ABCDF/1234/TTT`

	payload = `{"image" : "alpine"}`
)

func TestHandling(t *testing.T) {
	tests := []struct {
		caseDesc      string
		cfg           string
		expctdInvctns []invctn
	}{
		{
			"Single Route",
			singleRoute,
			[]invctn{
				{
					"*outputs.SlackOutput", "*regoservice.regoEvaluator", "route1",
				},
			},
		},
		{
			"2 Routes",
			twoRoutes,
			[]invctn{
				{
					"*outputs.SlackOutput", "*regoservice.regoEvaluator", "route1",
				},
				{
					"*outputs.SlackOutput", "*regoservice.regoEvaluator", "route2",
				},
			},
		},
		{
			"2 Outputs per single route",
			twoOutputs,
			[]invctn{
				{
					"*outputs.SlackOutput", "*regoservice.regoEvaluator", "route1",
				},
				{
					"*outputs.SlackOutput", "*regoservice.regoEvaluator", "route1",
				},
			},
		},
		{
			"No Outputs configured",
			noOutputs,
			[]invctn{},
		},
		{
			"No Template configured",
			noTemplates,
			[]invctn{},
		},
		{
			"Invalid Output reference",
			invalidOutput,
			[]invctn{},
		},
		{
			"Invalid Template reference",
			invalidTemplate,
			[]invctn{},
		},
		{
			"No outputs associated with route",
			noAssociatedOutput,
			[]invctn{},
		},
	}
	for _, test := range tests {
		runTestRouteHandlingCase(t, test.caseDesc, test.cfg, test.expctdInvctns)
	}
}
func runTestRouteHandlingCase(t *testing.T, caseDesc string, cfg string, expctdInvctns []invctn) {
	actualInvctCnt := 0
	t.Logf("Case: %s\n", caseDesc)
	wrap := ctxWrapper{}
	wrap.setup(cfg)

	defer wrap.teardown()

	err := wrap.instance.Start(wrap.cfgPath)
	if err != nil {
		t.Fatalf("[%s] Unexpected error %v", caseDesc, err)
	}

	wrap.instance.handle([]byte(payload))
	timeoutDuration := 3 * time.Second
	if len(expctdInvctns) == 0 {
		timeoutDuration = time.Second
	}
	timeout := time.After(timeoutDuration)
	for {
		select {
		case <-timeout:
			if len(expctdInvctns) > 0 {
				t.Fatal("test didn't finish in time")
			}
			return
		case r := <-wrap.buff:
			t.Logf("[%s] received invocation (%s, %s, %s)", caseDesc, r.routeName, r.outputCls, r.templateCls)
			actualInvctCnt++
			found := false
			for _, expect := range expctdInvctns {
				if r == expect {
					found = true
					break
				}
			}
			if actualInvctCnt == len(expctdInvctns) {
				return //everything is ok, exiting
			}
			if !found && len(expctdInvctns) > 0 {
				t.Errorf("[%s] Unexpected invocation (%s, %s, %s)", caseDesc, r.routeName, r.outputCls, r.templateCls)
				return
			}
			if actualInvctCnt > len(expctdInvctns) {
				t.Errorf("[%s] Service should be called %d times but called %d times", caseDesc, len(expctdInvctns), actualInvctCnt)
				return
			}
		}
	}

}
func TestInvalidRouteName(t *testing.T) {
	expctdInvctns := 0
	actualInvctCnt := 0
	wrap := ctxWrapper{}
	wrap.setup(singleRoute)

	defer wrap.teardown()

	err := wrap.instance.Start(wrap.cfgPath)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	wrap.instance.HandleRoute("not-exist", []byte(payload))
	timeout := time.After(1 * time.Second)
	for {
		select {
		case <-timeout:
			return
		case <-wrap.buff:
			actualInvctCnt++
			if actualInvctCnt > expctdInvctns {
				t.Errorf("Service shouldn't be called if invalid route is specified")
				return
			}
		}
	}

}
func TestSend(t *testing.T) {
	expctdInvctns := 1
	actualInvctCnt := 0
	wrap := ctxWrapper{}
	wrap.setup(singleRoute)

	defer wrap.teardown()

	err := wrap.instance.Start(wrap.cfgPath)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	wrap.instance.Send([]byte(payload))
	timeout := time.After(1 * time.Second)
	for {
		select {
		case <-timeout:
			return
		case <-wrap.buff:
			actualInvctCnt++
			if actualInvctCnt != expctdInvctns {
				t.Errorf("Service shouldn't be called once")
				return
			}
		}
	}
}
