"use strict";(self.webpackChunksiren=self.webpackChunksiren||[]).push([[690],{3905:function(e,t,n){n.d(t,{Zo:function(){return u},kt:function(){return d}});var r=n(7294);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function l(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},o=Object.keys(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var s=r.createContext({}),c=function(e){var t=r.useContext(s),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},u=function(e){var t=c(e.components);return r.createElement(s.Provider,{value:t},e.children)},p={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},m=r.forwardRef((function(e,t){var n=e.components,a=e.mdxType,o=e.originalType,s=e.parentName,u=l(e,["components","mdxType","originalType","parentName"]),m=c(n),d=a,g=m["".concat(s,".").concat(d)]||m[d]||p[d]||o;return n?r.createElement(g,i(i({ref:t},u),{},{components:n})):r.createElement(g,i({ref:t},u))}));function d(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var o=n.length,i=new Array(o);i[0]=m;var l={};for(var s in t)hasOwnProperty.call(t,s)&&(l[s]=t[s]);l.originalType=e,l.mdxType="string"==typeof e?e:a,i[1]=l;for(var c=2;c<o;c++)i[c]=n[c];return r.createElement.apply(null,i)}return r.createElement.apply(null,n)}m.displayName="MDXCreateElement"},7508:function(e,t,n){n.r(t),n.d(t,{frontMatter:function(){return l},contentTitle:function(){return s},metadata:function(){return c},toc:function(){return u},default:function(){return m}});var r=n(7462),a=n(3366),o=(n(7294),n(3905)),i=["components"],l={},s="Alert History Subscription",c={unversionedId:"guides/alert_history",id:"guides/alert_history",isDocsHomePage:!1,title:"Alert History Subscription",description:"Siren can store the alerts triggered via Cortex Alertmanager. Cortex alertmanager is configured to call Siren API, using",source:"@site/docs/guides/alert_history.md",sourceDirName:"guides",slug:"/guides/alert_history",permalink:"/siren/docs/guides/alert_history",editUrl:"https://github.com/odpf/siren/edit/master/docs/docs/guides/alert_history.md",tags:[],version:"current",lastUpdatedBy:"Rahmat Hidayat",lastUpdatedAt:1650440190,formattedLastUpdatedAt:"4/20/2022",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Templates",permalink:"/siren/docs/guides/templates"},next:{title:"Bulk Rule management",permalink:"/siren/docs/guides/bulk_rules"}},u=[],p={toc:u};function m(e){var t=e.components,l=(0,a.Z)(e,i);return(0,o.kt)("wrapper",(0,r.Z)({},p,l,{components:t,mdxType:"MDXLayout"}),(0,o.kt)("h1",{id:"alert-history-subscription"},"Alert History Subscription"),(0,o.kt)("p",null,"Siren can store the alerts triggered via Cortex Alertmanager. Cortex alertmanager is configured to call Siren API, using\na webhook receiver. This is done by adding a subscription using an HTTP Receiver on empty match condition, which will\nresult in calling thus HTTP Receiver on all alerts"),(0,o.kt)("p",null,"Example Receiver"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-json"},'{\n  "id": "1",\n  "name": "alert-history-receiver",\n  "type": "http",\n  "labels": {\n    "team": "siren-devs-alert-history"\n  },\n  "configurations": {\n    "url": "http://localhost:3000/v1beta1/alerts/cortex/3"\n  }\n}\n')),(0,o.kt)("p",null,"Note that the url has ",(0,o.kt)("inlineCode",{parentName:"p"},"cortex/3")," at the end, which means this will be able to parse alert history payloads from cortex\ntype and store in DB by making it belong to provider id ",(0,o.kt)("inlineCode",{parentName:"p"},"3"),"."),(0,o.kt)("p",null,"We will need the subscription as well, example:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-json"},'{\n  "id": "384",\n  "urn": "alert-history-subscription",\n  "namespace": "10",\n  "receivers": [\n    {\n      "id": "1"\n    }\n  ],\n  "match": {}\n}\n')),(0,o.kt)("p",null,"After this, as soon as any alert is sent by Alertmanager to slack or pagerduty, it will be sent to Siren for storage\npurpose."),(0,o.kt)("p",null,(0,o.kt)("img",{alt:"Siren Alert History",src:n(4324).Z})),(0,o.kt)("p",null,"The parsing of payload from alert manager depends on a particular syntax. you can configure your templates to follow\nthis syntax, with proper annotations to identify:"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},"Which alert was triggered"),(0,o.kt)("li",{parentName:"ul"},"Which resource this alert refers to"),(0,o.kt)("li",{parentName:"ul"},"On Which metric, this alert was triggered"),(0,o.kt)("li",{parentName:"ul"},"What was the metric value for alert trigger"),(0,o.kt)("li",{parentName:"ul"},"What was the severity of alert(CRITICAL, WARNING or RESOLVED)")),(0,o.kt)("p",null,"An Example template:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-yaml"},'apiVersion: v2\ntype: template\nname: CPU\nbody:\n  - alert: CPUWarning\n    expr: avg by (host) (cpu_usage_user{cpu="cpu-total"}) > [[.warning]]\n    for: "[[.for]]"\n    labels:\n      severity: WARNING\n    annotations:\n      description: CPU has been above [[.warning]] for last [[.for]] {{ $labels.host }}\n      resource: { { $labels.instance } }\n      template: CPU\n      metricName: cpu_usage_user\n      metricValue: { { $labels.cpu_usage_user } }\n  - alert: CPUCritical\n    expr: avg by (host) (cpu_usage_user{cpu="cpu-total"}) > [[.critical]]\n    for: "[[.for]]"\n    labels:\n      severity: CRITICAL\n    annotations:\n      description: CPU has been above [[.critical]] for last [[.for]] {{ $labels.host }}\n      resource: { { $labels.instance } }\n      template: CPU\n      metricName: cpu_usage_user\n      metricValue: { { $labels.cpu_usage_user } }\nvariables:\n  - name: for\n    type: string\n    default: 10m\n    description: For eg 5m, 2h; Golang duration format\n  - name: warning\n    type: int\n    default: 80\n  - name: critical\n    type: int\n    default: 90\ntags:\n  - systems\n')),(0,o.kt)("p",null,"Please note that, the mandatory keys, in order to successfully store Alert History is,"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-yaml"},"labels:\n  severity: CRITICAL\nannotations:\n  resource: { { $labels.instance } }\n  template: CPU\n  metricName: cpu_usage_user\n  metricValue: { { $labels.cpu_usage_user } }\n")),(0,o.kt)("p",null,"The keys are pretty obvious to match with what was described in bullets points in the introduction above."),(0,o.kt)("p",null,"In the above example we can see, the alert annotation has sufficient values for alert history storage. We can set up\ncortex alertmanager, to call Siren AlertHistory APIs as a webhook receiver. The above annotations and labels will be\nparsed by Siren APIs, to be stored in the database."),(0,o.kt)("p",null,(0,o.kt)("strong",{parentName:"p"},"Alert History Creation via API")),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-text"},'POST /v1beta1/alerts/cortex/1 HTTP/1.1\nHost: localhost:3000\nContent-Type: application/json\nContent-Length: 357\n\n{\n    "alerts": [\n        {\n            "status": "firing",\n            "labels": {\n                "severity": "CRITICAL"\n            },\n            "annotations": {\n                "resource": "apolloVM",\n                "template": "CPU",\n                "metricName": "cpu_usage_user",\n                "metricValue": "90"\n            }\n        }\n    ]\n}\n')),(0,o.kt)("p",null,"The request body of Alertmanager POST call to configured webhook looks something like (after you have followed the\nlabels and annotations c in the templates) above snippet."),(0,o.kt)("p",null,"The alerts API will parse the above payload and store in the database, which you can fetch via the GET APIs with proper\nfilters of startTime, endTime. See the swagger file for more details."))}m.isMDXComponent=!0},4324:function(e,t,n){t.Z=n.p+"assets/images/alerthistory-f454252342930a9749201a227cafac49.jpg"}}]);