"use strict";(self.webpackChunksiren=self.webpackChunksiren||[]).push([[774],{3905:function(e,t,n){n.d(t,{Zo:function(){return m},kt:function(){return c}});var a=n(7294);function r(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function l(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?l(Object(n),!0).forEach((function(t){r(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):l(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function o(e,t){if(null==e)return{};var n,a,r=function(e,t){if(null==e)return{};var n,a,r={},l=Object.keys(e);for(a=0;a<l.length;a++)n=l[a],t.indexOf(n)>=0||(r[n]=e[n]);return r}(e,t);if(Object.getOwnPropertySymbols){var l=Object.getOwnPropertySymbols(e);for(a=0;a<l.length;a++)n=l[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(r[n]=e[n])}return r}var p=a.createContext({}),s=function(e){var t=a.useContext(p),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},m=function(e){var t=s(e.components);return a.createElement(p.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return a.createElement(a.Fragment,{},t)}},d=a.forwardRef((function(e,t){var n=e.components,r=e.mdxType,l=e.originalType,p=e.parentName,m=o(e,["components","mdxType","originalType","parentName"]),d=s(n),c=r,g=d["".concat(p,".").concat(c)]||d[c]||u[c]||l;return n?a.createElement(g,i(i({ref:t},m),{},{components:n})):a.createElement(g,i({ref:t},m))}));function c(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var l=n.length,i=new Array(l);i[0]=d;var o={};for(var p in t)hasOwnProperty.call(t,p)&&(o[p]=t[p]);o.originalType=e,o.mdxType="string"==typeof e?e:r,i[1]=o;for(var s=2;s<l;s++)i[s]=n[s];return a.createElement.apply(null,i)}return a.createElement.apply(null,n)}d.displayName="MDXCreateElement"},9825:function(e,t,n){n.r(t),n.d(t,{frontMatter:function(){return o},contentTitle:function(){return p},metadata:function(){return s},toc:function(){return m},default:function(){return d}});var a=n(7462),r=n(3366),l=(n(7294),n(3905)),i=["components"],o={},p="Templates",s={unversionedId:"guides/templates",id:"guides/templates",isDocsHomePage:!1,title:"Templates",description:"Siren templates are an abstraction",source:"@site/docs/guides/templates.md",sourceDirName:"guides",slug:"/guides/templates",permalink:"/siren/docs/guides/templates",editUrl:"https://github.com/odpf/siren/edit/master/docs/docs/guides/templates.md",tags:[],version:"current",lastUpdatedBy:"Rahmat Hidayat",lastUpdatedAt:1654657110,formattedLastUpdatedAt:"6/8/2022",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Rules",permalink:"/siren/docs/guides/rules"},next:{title:"Alert History Subscription",permalink:"/siren/docs/guides/alert_history"}},m=[{value:"API interface",id:"api-interface",children:[{value:"Template creation/update",id:"template-creationupdate",children:[]},{value:"Terminology of the request body",id:"terminology-of-the-request-body",children:[]},{value:"Fetching a template",id:"fetching-a-template",children:[]},{value:"Deleting a template",id:"deleting-a-template",children:[]}]},{value:"CLI interface",id:"cli-interface",children:[{value:"Terminology",id:"terminology",children:[]}]}],u={toc:m};function d(e){var t=e.components,n=(0,r.Z)(e,i);return(0,l.kt)("wrapper",(0,a.Z)({},u,n,{components:t,mdxType:"MDXLayout"}),(0,l.kt)("h1",{id:"templates"},"Templates"),(0,l.kt)("p",null,"Siren templates are an abstraction\nover ",(0,l.kt)("a",{parentName:"p",href:"https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/"},"Prometheus rules"),". It\nutilises ",(0,l.kt)("a",{parentName:"p",href:"https://golang.org/pkg/text/template/"},"go-templates")," to provide implements data-driven templates for\ngenerating textual output. The template delimiter used is ",(0,l.kt)("inlineCode",{parentName:"p"},"[[")," and ",(0,l.kt)("inlineCode",{parentName:"p"},"]]"),"."),(0,l.kt)("p",null,"One can create templates using either HTTP APIs or CLI."),(0,l.kt)("h2",{id:"api-interface"},"API interface"),(0,l.kt)("h3",{id:"template-creationupdate"},"Template creation/update"),(0,l.kt)("p",null,"Templates can be created using Siren APIs. The below snippet describes an example."),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre",className:"language-text"},'PUT /v1beta1/templates HTTP/1.1\nHost: localhost:3000\nContent-Type: application/json\nContent-Length: 1383\n\n{\n    "name": "CPU",\n    "body": "- alert: CPUHighWarning\\n  expr: avg by (host) (cpu_usage_user{cpu=\\"cpu-total\\"}) > [[.warning]]\\n  for: \'[[.for]]\'\\n  labels:\\n    severity: WARNING\\n    team: \'[[ .team ]]\'\\n  annotations:\\n    dashboard: https://example.com\\n    description: CPU has been above [[.warning]] for last [[.for]] {{ $labels.host }}\\n- alert: CPUHighCritical\\n  expr: avg by (host) (cpu_usage_user{cpu=\\"cpu-total\\"}) > [[.critical]]\\n  for: \'[[.for]]\'\\n  labels:\\n    severity: CRITICAL\\n    team: \'[[ .team ]]\'\\n  annotations:\\n    dashboard: example.com\\n    description: CPU has been above [[.critical]] for last [[.for]] {{ $labels.host }}\\n",\n    "tags": [\n        "firehose",\n        "dagger"\n    ],\n    "variables": [\n        {\n            "name": "team",\n            "type": "string",\n            "default": "odpf",\n            "description": "Name of the team that owns the deployment"\n        },\n        {\n            "name": "for",\n            "type": "string",\n            "default": "10m",\n            "description": "For eg 5m, 2h; Golang duration format"\n        },\n        {\n            "name": "warning",\n            "type": "int",\n            "default": "85",\n            "description": ""\n        },\n        {\n            "name": "critical",\n            "type": "int",\n            "default": "95",\n            "description": ""\n        }\n    ]\n}\n\n')),(0,l.kt)("h3",{id:"terminology-of-the-request-body"},"Terminology of the request body"),(0,l.kt)("table",null,(0,l.kt)("thead",{parentName:"table"},(0,l.kt)("tr",{parentName:"thead"},(0,l.kt)("th",{parentName:"tr",align:null},"Term"),(0,l.kt)("th",{parentName:"tr",align:null},"Description"),(0,l.kt)("th",{parentName:"tr",align:null},"Example/Default"))),(0,l.kt)("tbody",{parentName:"table"},(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},"Name"),(0,l.kt)("td",{parentName:"tr",align:null},"Name of the template"),(0,l.kt)("td",{parentName:"tr",align:null},"CPUHigh")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},"Body"),(0,l.kt)("td",{parentName:"tr",align:null},"Array of rule body. The body can be templatized in go template format."),(0,l.kt)("td",{parentName:"tr",align:null},"See example above")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},"Variables"),(0,l.kt)("td",{parentName:"tr",align:null},"Array of variables that were templatized in the body with their data type, default value and description."),(0,l.kt)("td",{parentName:"tr",align:null},"See example above")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},"Tags"),(0,l.kt)("td",{parentName:"tr",align:null},"Array of resources/applications that can utilize this template"),(0,l.kt)("td",{parentName:"tr",align:null},"VM")))),(0,l.kt)("p",null,"The response body will look like this:"),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre",className:"language-json"},'{\n  "id": 38,\n  "CreatedAt": "2021-04-29T16:20:48.061862+05:30",\n  "UpdatedAt": "2021-04-29T16:22:19.978837+05:30",\n  "name": "CPU",\n  "body": "- alert: CPUHighWarning\\n  expr: avg by (host) (cpu_usage_user{cpu=\\"cpu-total\\"}) > [[.warning]]\\n  for: \'[[.for]]\'\\n  labels:\\n    severity: WARNING\\n    team: \'[[ .team ]]\'\\n  annotations:\\n    dashboard: https://example.com\\n    description: CPU has been above [[.warning]] for last [[.for]] {{ $labels.host }}\\n- alert: CPUHighCritical\\n  expr: avg by (host) (cpu_usage_user{cpu=\\"cpu-total\\"}) > [[.critical]]\\n  for: \'[[.for]]\'\\n  labels:\\n    severity: CRITICAL\\n    team: \'[[ .team ]]\'\\n  annotations:\\n    dashboard: example.com\\n    description: CPU has been above [[.critical]] for last [[.for]] {{ $labels.host }}\\n",\n  "tags": [\n    "firehose",\n    "dagger"\n  ],\n  "variables": [\n    {\n      "name": "team",\n      "type": "string",\n      "default": "odpf",\n      "description": "Name of the team that owns the deployment"\n    },\n    {\n      "name": "for",\n      "type": "string",\n      "default": "10m",\n      "description": "For eg 5m, 2h; Golang duration format"\n    },\n    {\n      "name": "warning",\n      "type": "int",\n      "default": "85",\n      "description": ""\n    },\n    {\n      "name": "critical",\n      "type": "int",\n      "default": "95",\n      "description": ""\n    }\n  ]\n}\n')),(0,l.kt)("h3",{id:"fetching-a-template"},"Fetching a template"),(0,l.kt)("p",null,(0,l.kt)("strong",{parentName:"p"},"Fetching by Name")),(0,l.kt)("p",null,"Here is an example to fetch a template using name."),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre",className:"language-text"},"GET /v1beta1/templates/cpu HTTP/1.1\nHost: localhost:3000\n")),(0,l.kt)("p",null,(0,l.kt)("strong",{parentName:"p"},"Fetching by Tags")),(0,l.kt)("p",null,"Here is an example to fetch a templates matching the tag."),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre",className:"language-text"},"GET /v1beta1/templates?tag=firehose HTTP/1.1\nHost: localhost:3000\n")),(0,l.kt)("h3",{id:"deleting-a-template"},"Deleting a template"),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre",className:"language-text"},"DELETE /v1beta1/templates/cpu HTTP/1.1\nHost: localhost:3000\n")),(0,l.kt)("p",null,(0,l.kt)("strong",{parentName:"p"},"Note:")),(0,l.kt)("ol",null,(0,l.kt)("li",{parentName:"ol"},"Updating a template via API will not upload the associated rules.")),(0,l.kt)("h2",{id:"cli-interface"},"CLI interface"),(0,l.kt)("p",null,"With CLI, you will need a YAML file in the below specified format to create/update templates. The CLI calls Siren\nservice templates APIs in turn."),(0,l.kt)("p",null,(0,l.kt)("strong",{parentName:"p"},"Example template file")),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre",className:"language-yaml"},'apiVersion: v2\ntype: template\nname: CPU\nbody:\n  - alert: CPUWarning\n    expr: avg by (host) (cpu_usage_user{cpu="cpu-total"}) > [[.warning]]\n    for: "[[.for]]"\n    labels:\n      severity: WARNING\n    annotations:\n      description: CPU has been above [[.warning]] for last [[.for]] {{ $labels.host }}\n  - alert: CPUCritical\n    expr: avg by (host) (cpu_usage_user{cpu="cpu-total"}) > [[.critical]]\n    for: "[[.for]]"\n    labels:\n      severity: CRITICAL\n    annotations:\n      description: CPU has been above [[.critical]] for last [[.for]] {{ $labels.host }}\nvariables:\n  - name: for\n    type: string\n    default: 10m\n    description: For eg 5m, 2h; Golang duration format\n  - name: warning\n    type: int\n    default: 80\n  - name: critical\n    type: int\n    default: 90\ntags:\n  - systems\n')),(0,l.kt)("p",null,"In the above example, we are using one template to define rules of two severity labels viz WARNING and CRITICAL. Here we\nhave made 3 templates variables ",(0,l.kt)("inlineCode",{parentName:"p"},"for"),", ",(0,l.kt)("inlineCode",{parentName:"p"},"warning")," and ",(0,l.kt)("inlineCode",{parentName:"p"},"critical")," which denote the appropriate alerting thresholds. They\nwill be given a value while actual rule(alert) creating."),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre",className:"language-shell"},"go run main.go upload cpu_template.yaml\n")),(0,l.kt)("h3",{id:"terminology"},"Terminology"),(0,l.kt)("table",null,(0,l.kt)("thead",{parentName:"table"},(0,l.kt)("tr",{parentName:"thead"},(0,l.kt)("th",{parentName:"tr",align:null},"Term"),(0,l.kt)("th",{parentName:"tr",align:null},"Description"),(0,l.kt)("th",{parentName:"tr",align:null},"Example/Default"))),(0,l.kt)("tbody",{parentName:"table"},(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},"API Version"),(0,l.kt)("td",{parentName:"tr",align:null},"Which API to use to parse the YAML file"),(0,l.kt)("td",{parentName:"tr",align:null},"v2")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},"Type"),(0,l.kt)("td",{parentName:"tr",align:null},"Describes the type of object represented by YAML file"),(0,l.kt)("td",{parentName:"tr",align:null},"template")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},"Name"),(0,l.kt)("td",{parentName:"tr",align:null},"Name of the template"),(0,l.kt)("td",{parentName:"tr",align:null},"CPUHigh")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},"Body"),(0,l.kt)("td",{parentName:"tr",align:null},"Array of rule body. The body can be templatized in go template format."),(0,l.kt)("td",{parentName:"tr",align:null},"See example file")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},"Variables"),(0,l.kt)("td",{parentName:"tr",align:null},"Array of variables that were templatized in the body with their data type, default value and description."),(0,l.kt)("td",{parentName:"tr",align:null},"See example file")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},"Tags"),(0,l.kt)("td",{parentName:"tr",align:null},"Array of resources/applications that can utilize this template"),(0,l.kt)("td",{parentName:"tr",align:null},"VM")))),(0,l.kt)("p",null,(0,l.kt)("strong",{parentName:"p"},"Note:")),(0,l.kt)("ol",null,(0,l.kt)("li",{parentName:"ol"},"It's suggested to always provide default value for the templated variables."),(0,l.kt)("li",{parentName:"ol"},"Updating a template via CLI will update all associated rules.")))}d.isMDXComponent=!0}}]);