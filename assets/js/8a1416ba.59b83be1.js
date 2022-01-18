"use strict";(self.webpackChunksiren=self.webpackChunksiren||[]).push([[886],{3905:function(e,t,i){i.d(t,{Zo:function(){return c},kt:function(){return o}});var n=i(7294);function r(e,t,i){return t in e?Object.defineProperty(e,t,{value:i,enumerable:!0,configurable:!0,writable:!0}):e[t]=i,e}function a(e,t){var i=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),i.push.apply(i,n)}return i}function M(e){for(var t=1;t<arguments.length;t++){var i=null!=arguments[t]?arguments[t]:{};t%2?a(Object(i),!0).forEach((function(t){r(e,t,i[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(i)):a(Object(i)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(i,t))}))}return e}function I(e,t){if(null==e)return{};var i,n,r=function(e,t){if(null==e)return{};var i,n,r={},a=Object.keys(e);for(n=0;n<a.length;n++)i=a[n],t.indexOf(i)>=0||(r[i]=e[i]);return r}(e,t);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(n=0;n<a.length;n++)i=a[n],t.indexOf(i)>=0||Object.prototype.propertyIsEnumerable.call(e,i)&&(r[i]=e[i])}return r}var g=n.createContext({}),l=function(e){var t=n.useContext(g),i=t;return e&&(i="function"==typeof e?e(t):M(M({},t),e)),i},c=function(e){var t=l(e.components);return n.createElement(g.Provider,{value:t},e.children)},N={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},u=n.forwardRef((function(e,t){var i=e.components,r=e.mdxType,a=e.originalType,g=e.parentName,c=I(e,["components","mdxType","originalType","parentName"]),u=l(i),o=r,s=u["".concat(g,".").concat(o)]||u[o]||N[o]||a;return i?n.createElement(s,M(M({ref:t},c),{},{components:i})):n.createElement(s,M({ref:t},c))}));function o(e,t){var i=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var a=i.length,M=new Array(a);M[0]=u;var I={};for(var g in t)hasOwnProperty.call(t,g)&&(I[g]=t[g]);I.originalType=e,I.mdxType="string"==typeof e?e:r,M[1]=I;for(var l=2;l<a;l++)M[l]=i[l];return n.createElement.apply(null,M)}return n.createElement.apply(null,i)}u.displayName="MDXCreateElement"},4730:function(e,t,i){i.r(t),i.d(t,{frontMatter:function(){return I},contentTitle:function(){return g},metadata:function(){return l},toc:function(){return c},default:function(){return u}});var n=i(7462),r=i(3366),a=(i(7294),i(3905)),M=["components"],I={},g="Architecture",l={unversionedId:"concepts/architecture",id:"concepts/architecture",isDocsHomePage:!1,title:"Architecture",description:"Siren exposes HTTP API to allow rule, template and slack & Pagerduty credentials configuration. It talks to upstream",source:"@site/docs/concepts/architecture.md",sourceDirName:"concepts",slug:"/concepts/architecture",permalink:"/siren/docs/concepts/architecture",editUrl:"https://github.com/odpf/siren/edit/master/docs/docs/concepts/architecture.md",tags:[],version:"current",lastUpdatedBy:"Ravi Suhag",lastUpdatedAt:1642494349,formattedLastUpdatedAt:"1/18/2022",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Overview",permalink:"/siren/docs/concepts/overview"},next:{title:"Schema Design",permalink:"/siren/docs/concepts/schema"}},c=[{value:"System Design",id:"system-design",children:[{value:"Technologies",id:"technologies",children:[]},{value:"Components",id:"components",children:[]},{value:"Managing Templates via YAML File",id:"managing-templates-via-yaml-file",children:[]},{value:"Managing Rules via YAML File",id:"managing-rules-via-yaml-file",children:[]}]},{value:"Siren Integration",id:"siren-integration",children:[{value:"Cortex Ruler",id:"cortex-ruler",children:[]},{value:"Cortex Alertmanager",id:"cortex-alertmanager",children:[]}]}],N={toc:c};function u(e){var t=e.components,I=(0,r.Z)(e,M);return(0,a.kt)("wrapper",(0,n.Z)({},N,I,{components:t,mdxType:"MDXLayout"}),(0,a.kt)("h1",{id:"architecture"},"Architecture"),(0,a.kt)("p",null,"Siren exposes HTTP API to allow rule, template and slack & Pagerduty credentials configuration. It talks to upstream\ncortex ruler to configure rules(alerting and recording rules). It talks to Cortex Alertmanager to configure the\ndestination where alerts should go. It stores data around credentials, templates and current state of configured alerts\nin PostgresDB. It also stores alerts triggered via Cortex Alertmanager."),(0,a.kt)("p",null,(0,a.kt)("img",{alt:"Siren Architecture",src:i(6044).Z})),(0,a.kt)("h2",{id:"system-design"},"System Design"),(0,a.kt)("h3",{id:"technologies"},"Technologies"),(0,a.kt)("p",null,"Siren is developed with"),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},"Golang - Programming language"),(0,a.kt)("li",{parentName:"ul"},"Docker - container engine to start postgres and cortex to aid development"),(0,a.kt)("li",{parentName:"ul"},"Cortex - multi-tenant prometheus based monitoring stack"),(0,a.kt)("li",{parentName:"ul"},"Postgres - a relational database")),(0,a.kt)("h3",{id:"components"},"Components"),(0,a.kt)("p",null,(0,a.kt)("em",{parentName:"p"},(0,a.kt)("strong",{parentName:"em"},"GRPC Server and HTTP Gateway"))),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},"GRPC Server exposes RPC APIs and RESTfull APIs (via GRPC gateway) to allow configuration of rules, templates, alerting\ncredentials and storing triggered alert history.")),(0,a.kt)("p",null,(0,a.kt)("em",{parentName:"p"},(0,a.kt)("strong",{parentName:"em"},"PostgresDB"))),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},"Used for storing the templates in a predefined schema enabling reuse of same rule body."),(0,a.kt)("li",{parentName:"ul"},"Stores the rules configured via HTTP APIs and used for preserving thresholds when rule is deleted"),(0,a.kt)("li",{parentName:"ul"},"Stores Slack and Pagerduty credentials to enable DIY interface for configuring destinations for alerting.")),(0,a.kt)("p",null,(0,a.kt)("em",{parentName:"p"},(0,a.kt)("strong",{parentName:"em"},"Command Line Interface"))),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},"Provides a way to manage rules and templates using YAML files in a format described below."),(0,a.kt)("li",{parentName:"ul"},"Run a web-server from CLI."),(0,a.kt)("li",{parentName:"ul"},"Runs DB Migrations"),(0,a.kt)("li",{parentName:"ul"},"Manage templates, rules, providers, namespaces and receivers")),(0,a.kt)("h3",{id:"managing-templates-via-yaml-file"},"Managing Templates via YAML File"),(0,a.kt)("p",null,"Siren gives flexibility to templatize prometheus rules for re-usability purpose. Template can be managed via APIs(REST\nand GRPC). Apart from that, there is a command line interface as well which parses a YAML file in a specified format (as\ndescribed below) and upload to Siren using an HTTP Client of Siren Service. Refer ",(0,a.kt)("a",{parentName:"p",href:"/siren/docs/guides/templates"},"here")," for\nmore details around usage and terminology."),(0,a.kt)("h3",{id:"managing-rules-via-yaml-file"},"Managing Rules via YAML File"),(0,a.kt)("p",null,"To manage rules in bulk, Siren gives a way to manage rules using YAML files, which you can manage in a version\ncontrolled repository. Using the ",(0,a.kt)("inlineCode",{parentName:"p"},"upload")," command one can upload a rule YAML file in a specified format (as described\nbelow) and upload to Siren using the GRPC Client(comes inbuilt) of Siren Service. Refer ",(0,a.kt)("a",{parentName:"p",href:"/siren/docs/guides/rules"},"here")," for\nmore details around usage and terminology."),(0,a.kt)("p",null,(0,a.kt)("strong",{parentName:"p"},"Note:")," Updating a template also updates the associated rules."),(0,a.kt)("h2",{id:"siren-integration"},"Siren Integration"),(0,a.kt)("p",null,"The monitoring providers supported are:"),(0,a.kt)("ol",null,(0,a.kt)("li",{parentName:"ol"},"Cortex metrics.")),(0,a.kt)("p",null,"The section details all integrating systems for Siren deployment."),(0,a.kt)("h3",{id:"cortex-ruler"},"Cortex Ruler"),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},"The upstream Cortex ruler deployment which is used for rule creation in the proper namespace/group. You can create\na ",(0,a.kt)("a",{parentName:"li",href:"/siren/docs/guides/providers"},"provider")," for that purpose and provide appropriate hostname.")),(0,a.kt)("h3",{id:"cortex-alertmanager"},"Cortex Alertmanager"),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},"The upstream Cortex alertmanager deployment where routing configurations are stored in the proper format. Sirenstores\nsubscriptions which gets synced in the alertmanager. Cortex Alertmanger hostname is fetched\nfrom ",(0,a.kt)("a",{parentName:"li",href:"/siren/docs/guides/providers"},"provider's")," host key.")))}u.isMDXComponent=!0},6044:function(e,t){t.Z="data:image/svg+xml;base64,PHN2ZyB2ZXJzaW9uPSIxLjEiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgdmlld0JveD0iMCAwIDc0OC40NzAwMzE3MzgyODEyIDI2Ny4yNjIxNzY1MTM2NzE5IiB3aWR0aD0iMjI0NS40MTAwOTUyMTQ4NDM4IiBoZWlnaHQ9IjgwMS43ODY1Mjk1NDEwMTU2Ij4KICA8IS0tIHN2Zy1zb3VyY2U6ZXhjYWxpZHJhdyAtLT4KICAKICA8ZGVmcz4KICAgIDxzdHlsZT4KICAgICAgQGZvbnQtZmFjZSB7CiAgICAgICAgZm9udC1mYW1pbHk6ICJWaXJnaWwiOwogICAgICAgIHNyYzogdXJsKCJodHRwczovL2V4Y2FsaWRyYXcuY29tL1ZpcmdpbC53b2ZmMiIpOwogICAgICB9CiAgICAgIEBmb250LWZhY2UgewogICAgICAgIGZvbnQtZmFtaWx5OiAiQ2FzY2FkaWEiOwogICAgICAgIHNyYzogdXJsKCJodHRwczovL2V4Y2FsaWRyYXcuY29tL0Nhc2NhZGlhLndvZmYyIik7CiAgICAgIH0KICAgIDwvc3R5bGU+CiAgPC9kZWZzPgogIDxyZWN0IHg9IjAiIHk9IjAiIHdpZHRoPSI3NDguNDcwMDMxNzM4MjgxMiIgaGVpZ2h0PSIyNjcuMjYyMTc2NTEzNjcxOSIgZmlsbD0iI2ZmZmZmZiI+PC9yZWN0PjxnIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIgdHJhbnNmb3JtPSJ0cmFuc2xhdGUoMTMuMTI4OTA2MjUgMjAuMjY0MTI5NjM4NjcxODc1KSByb3RhdGUoMCA2OS4wNTUxNDUyNjM2NzE4OCAyOC4wMTk4MjExNjY5OTIxODgpIj48cGF0aCBkPSJNMCAtMC44OCBDMzQuMDQgMS43NSwgNjcuMDQgMS41LCAxMzguMDMgLTEuMyBNLTAuNCAwLjgyIEMyNy41OCAxLjY4LCA1Ny42OCAxLjQsIDEzNy40NCAtMC4zNyBNMTM3LjcyIDEuNDYgQzEzOC40IDEzLjc1LCAxMzguODcgMjUuMzEsIDEzOC4wMyA1Ni41NSBNMTM3LjQ5IDAuNDUgQzEzOC42MSAxNS4zNiwgMTM2Ljg4IDMxLjExLCAxMzcuNTEgNTYuNTYgTTEzNy40NSA1Ni45OCBDMTA3Ljk5IDU1LjI5LCA3My4xMSA1Ni4yNSwgMCA1NS43MSBNMTM4LjQzIDU1LjkgQzkwLjA2IDU1Ljc3LCA0Mi42NCA1NC42MywgLTAuODcgNTUuOTkgTS0xLjg5IDU2LjQ4IEMtMC4yNSAzOS4yNSwgMS43NiAyMi4zOSwgMS4zMSAtMC4yOCBNLTAuMzEgNTYuNzkgQy0wLjE2IDM2LjAxLCAtMC4wOCAxNy4zMiwgLTAuMSAtMC4yIiBzdHJva2U9IiMwMDAwMDAiIHN0cm9rZS13aWR0aD0iMSIgZmlsbD0ibm9uZSI+PC9wYXRoPjwvZz48ZyB0cmFuc2Zvcm09InRyYW5zbGF0ZSg1NC44MTI2MjIwNzAzMTI1IDM4LjAzNDkxMjEwOTM3NSkgcm90YXRlKDAgMjggMTIuNSkiPjx0ZXh0IHg9IjAiIHk9IjE4IiBmb250LWZhbWlseT0iVmlyZ2lsLCBTZWdvZSBVSSBFbW9qaSIgZm9udC1zaXplPSIyMHB4IiBmaWxsPSIjMDAwMDAwIiB0ZXh0LWFuY2hvcj0ic3RhcnQiIHN0eWxlPSJ3aGl0ZS1zcGFjZTogcHJlOyIgZGlyZWN0aW9uPSJsdHIiPlNpcmVuIDwvdGV4dD48L2c+PGcgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIiB0cmFuc2Zvcm09InRyYW5zbGF0ZSg1MTcuNDcwMDMxNzM4MjgxMiAxNi42MjY0NjQ4NDM3NSkgcm90YXRlKDAgMTEwLjUgMjguMTgxNTY0MzMxMDU0Njg4KSI+PHBhdGggZD0iTTAuOTQgLTAuNjUgQzU4LjYgMC41OCwgMTE4LjY1IDIuNjUsIDIyMi4xNiAxLjE0IE0tMC42IDAuNTUgQzY2LjU0IC0wLjI0LCAxMzQuNDQgLTAuNDgsIDIyMS44NiAwLjcgTTIyMC43MSAtMC44OCBDMjIxLjI3IDE4Ljk5LCAyMjAuOTIgMzcuODIsIDIyMS42NSA1Ni4xNyBNMjIwLjk0IDAuMTYgQzIyMS43MiAyMS42MiwgMjIwLjQ1IDQzLjg1LCAyMjAuNjYgNTYuMjEgTTIyMC45NyA1Ni40NyBDMTczLjA0IDU3Ljg2LCAxMjIuNDMgNTguMDMsIDEuNzIgNTYuMjUgTTIyMS43OSA1NS41NyBDMTUzLjEyIDU1LjQyLCA4Mi43NCA1NS4wMSwgMC4xOCA1Ni40OSBNLTEuNzYgNTcuMjQgQzAuODIgMzUuMDksIC0xLjgyIDE0LjEsIDEuODkgMC45MiBNMC43OSA1Ni4xMyBDLTAuMzcgMzUuNDIsIDAuMjkgMTMuNjEsIC0wLjYxIC0wLjgyIiBzdHJva2U9IiMwMDAwMDAiIHN0cm9rZS13aWR0aD0iMSIgZmlsbD0ibm9uZSI+PC9wYXRoPjwvZz48ZyB0cmFuc2Zvcm09InRyYW5zbGF0ZSg1MjIuNDcwMDMxNzM4MjgxMiAzMi4zMDgwMjkxNzQ4MDQ2OSkgcm90YXRlKDAgMTA1LjUgMTIuNSkiPjx0ZXh0IHg9IjEwNS41IiB5PSIxOCIgZm9udC1mYW1pbHk9IlZpcmdpbCwgU2Vnb2UgVUkgRW1vamkiIGZvbnQtc2l6ZT0iMjBweCIgZmlsbD0iIzAwMDAwMCIgdGV4dC1hbmNob3I9Im1pZGRsZSIgc3R5bGU9IndoaXRlLXNwYWNlOiBwcmU7IiBkaXJlY3Rpb249Imx0ciI+TW9uaXRvcmluZyBQcm92aWRlcjwvdGV4dD48L2c+PGcgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIj48ZyB0cmFuc2Zvcm09InRyYW5zbGF0ZSgxNTMuNzEwNzU0Mzk0NTMxMjUgMzIuMjk2Nzk1ODg0MzY1NjgpIHJvdGF0ZSgwIDE4MS4zMjA3NDQyNjQ2Mjg3MyAyLjIxNDM4MDYyNjY0NDgxMSkiPjxwYXRoIGQ9Ik0wLjMyIDAuODYgQzYxLjA4IDEuNDksIDMwMy4yNiAyLjQ4LCAzNjMuNjIgMy4yMSBNLTAuOTggMC4yNyBDNTkuNzkgMS4wNSwgMzAyLjcgMy42NSwgMzYzLjEyIDQuMTYiIHN0cm9rZT0iIzAwMDAwMCIgc3Ryb2tlLXdpZHRoPSIxIiBmaWxsPSJub25lIj48L3BhdGg+PC9nPjxnIHRyYW5zZm9ybT0idHJhbnNsYXRlKDE1My43MTA3NTQzOTQ1MzEyNSAzMi4yOTY3OTU4ODQzNjU2OCkgcm90YXRlKDAgMTgxLjMyMDc0NDI2NDYyODczIDIuMjE0MzgwNjI2NjQ0ODExKSI+PHBhdGggZD0iTTMzNi4yNyAxNS41NyBDMzM5LjE4IDExLjQsIDM0OC42NyA5LjI1LCAzNjMuNjUgNS4wMyBNMzM1LjI2IDEzLjg5IEMzNDMuNiAxMS4xOCwgMzUwLjk2IDcuMjUsIDM2Mi44MiA0LjQ1IiBzdHJva2U9IiMwMDAwMDAiIHN0cm9rZS13aWR0aD0iMSIgZmlsbD0ibm9uZSI+PC9wYXRoPjwvZz48ZyB0cmFuc2Zvcm09InRyYW5zbGF0ZSgxNTMuNzEwNzU0Mzk0NTMxMjUgMzIuMjk2Nzk1ODg0MzY1NjgpIHJvdGF0ZSgwIDE4MS4zMjA3NDQyNjQ2Mjg3MyAyLjIxNDM4MDYyNjY0NDgxMSkiPjxwYXRoIGQ9Ik0zMzYuNDcgLTQuOTUgQzMzOS4yNCAtNC42NiwgMzQ4LjY5IC0yLjM0LCAzNjMuNjUgNS4wMyBNMzM1LjQ2IC02LjYzIEMzNDMuNzEgLTMuMjMsIDM1MSAtMS4wNiwgMzYyLjgyIDQuNDUiIHN0cm9rZT0iIzAwMDAwMCIgc3Ryb2tlLXdpZHRoPSIxIiBmaWxsPSJub25lIj48L3BhdGg+PC9nPjwvZz48ZyB0cmFuc2Zvcm09InRyYW5zbGF0ZSgyMjYuNzA3MzM2NDI1NzgxMjUgMTApIHJvdGF0ZSgwIDEyNSAyNSkiPjx0ZXh0IHg9IjAiIHk9IjE4IiBmb250LWZhbWlseT0iVmlyZ2lsLCBTZWdvZSBVSSBFbW9qaSIgZm9udC1zaXplPSIyMHB4IiBmaWxsPSIjMDAwMDAwIiB0ZXh0LWFuY2hvcj0ic3RhcnQiIHN0eWxlPSJ3aGl0ZS1zcGFjZTogcHJlOyIgZGlyZWN0aW9uPSJsdHIiPlJ1bGUgbWFuYWdlbWVudDwvdGV4dD48dGV4dCB4PSIwIiB5PSI0MyIgZm9udC1mYW1pbHk9IlZpcmdpbCwgU2Vnb2UgVUkgRW1vamkiIGZvbnQtc2l6ZT0iMjBweCIgZmlsbD0iIzAwMDAwMCIgdGV4dC1hbmNob3I9InN0YXJ0IiBzdHlsZT0id2hpdGUtc3BhY2U6IHByZTsiIGRpcmVjdGlvbj0ibHRyIj5BbGVydCByb3V0aW5nIG1hbmFnZW1lbnQ8L3RleHQ+PC9nPjxnIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIgdHJhbnNmb3JtPSJ0cmFuc2xhdGUoMTAgMjEyLjQ4MzMwNjg4NDc2NTYyKSByb3RhdGUoMCA2OCAyMi4zODk0MzQ4MTQ0NTMxMjUpIj48cGF0aCBkPSJNMC40IC0xLjAxIEM1MS43NyAxLjc2LCAxMDUuNDQgLTAuOTgsIDEzNS4zOCAxLjg5IE0tMC4wNCAwLjE2IEM0OS4xNyAtMC41LCA5OS4zNSAtMC42MiwgMTM2LjY2IDAuMzUgTTEzNC40NSAwLjM3IEMxMzcuMzMgMTEuOTgsIDEzNi4wMyAyMi42NywgMTM2LjAyIDQ0LjUgTTEzNS4wNSAtMC4wNiBDMTM2LjkzIDEyLjA0LCAxMzcuMzUgMjMuODcsIDEzNi45MSA0NC42MSBNMTM3LjcxIDQ1LjQxIEM5MS43MyA0My44NSwgNDQuOTcgNDMuMjYsIC0xLjE1IDQzLjc4IE0xMzYuMDkgNDQuNzIgQzEwNC4yNCA0NC4yNSwgNzQuNjkgNDMuOSwgMC42NiA0NS4zMyBNLTAuOCA0NS43MiBDLTEuODIgMzIuMDEsIC0xLjM0IDI0LjQzLCAtMS4xMiAxLjggTTAuMzUgNDUuNDUgQy0wLjQ4IDMyLjI2LCAtMC4xNiAxNy45NSwgLTAuNTcgLTAuOTIiIHN0cm9rZT0iIzAwMDAwMCIgc3Ryb2tlLXdpZHRoPSIxIiBmaWxsPSJub25lIj48L3BhdGg+PC9nPjxnIHRyYW5zZm9ybT0idHJhbnNsYXRlKDE0Ljc2NjQ5OTk5OTk5OTk1IDIyMi4zNzI3NDE2OTkyMTg3NSkgcm90YXRlKDAgNjMgMTIuNSkiPjx0ZXh0IHg9IjYzIiB5PSIxOCIgZm9udC1mYW1pbHk9IlZpcmdpbCwgU2Vnb2UgVUkgRW1vamkiIGZvbnQtc2l6ZT0iMjBweCIgZmlsbD0iIzAwMDAwMCIgdGV4dC1hbmNob3I9Im1pZGRsZSIgc3R5bGU9IndoaXRlLXNwYWNlOiBwcmU7IiBkaXJlY3Rpb249Imx0ciI+UG9zdGdyZXNEQjwvdGV4dD48L2c+PGcgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIj48ZyB0cmFuc2Zvcm09InRyYW5zbGF0ZSg3MC41NjkzNjk1MDM3MjE2NSA3OS41NDA1MjA4MzE0NTgyNSkgcm90YXRlKDAgLTAuNjAxMjQ2MTk0NzE2NTY2MiA2NS43NzQ4Njk2NTU0MjY1OCkiPjxwYXRoIGQ9Ik0wLjM0IDAuMSBDMC4yNyAyMi4wMSwgLTAuNTUgMTA5LjMsIC0wLjkgMTMxLjI0IE0tMC45NCAtMC45IEMtMS4wOCAyMS4xOSwgLTEuMjUgMTEwLjEzLCAtMS41NCAxMzIuNDUiIHN0cm9rZT0iIzAwMDAwMCIgc3Ryb2tlLXdpZHRoPSIxIiBmaWxsPSJub25lIj48L3BhdGg+PC9nPjxnIHRyYW5zZm9ybT0idHJhbnNsYXRlKDcwLjU2OTM2OTUwMzcyMTY1IDc5LjU0MDUyMDgzMTQ1ODI1KSByb3RhdGUoMCAtMC42MDEyNDYxOTQ3MTY1NjYyIDY1Ljc3NDg2OTY1NTQyNjU4KSI+PHBhdGggZD0iTS0xMS40NSAxMDUuMTYgQy0xMC4xIDExMS4wOSwgLTYuNjQgMTE5LjcsIC0wLjM0IDEzMy4yNCBNLTEyLjExIDEwMy42NSBDLTguMzMgMTE1LjMzLCAtMy44NSAxMjUuNDQsIC0xLjIyIDEzMS45NiIgc3Ryb2tlPSIjMDAwMDAwIiBzdHJva2Utd2lkdGg9IjEiIGZpbGw9Im5vbmUiPjwvcGF0aD48L2c+PGcgdHJhbnNmb3JtPSJ0cmFuc2xhdGUoNzAuNTY5MzY5NTAzNzIxNjUgNzkuNTQwNTIwODMxNDU4MjUpIHJvdGF0ZSgwIC0wLjYwMTI0NjE5NDcxNjU2NjIgNjUuNzc0ODY5NjU1NDI2NTgpIj48cGF0aCBkPSJNOS4wNyAxMDUuMyBDNS4wNSAxMTEuMjUsIDMuMTQgMTE5LjgzLCAtMC4zNCAxMzMuMjQgTTguNDEgMTAzLjc4IEM0LjYyIDExNS40OSwgMS41MyAxMjUuNTYsIC0xLjIyIDEzMS45NiIgc3Ryb2tlPSIjMDAwMDAwIiBzdHJva2Utd2lkdGg9IjEiIGZpbGw9Im5vbmUiPjwvcGF0aD48L2c+PC9nPjxnIHRyYW5zZm9ybT0idHJhbnNsYXRlKDc1LjMwNzczOTI1NzgxMjUgOTQuODM2MjczMTkzMzU5MzgpIHJvdGF0ZSgwIDczLjUgMzcuNSkiPjx0ZXh0IHg9IjAiIHk9IjE4IiBmb250LWZhbWlseT0iVmlyZ2lsLCBTZWdvZSBVSSBFbW9qaSIgZm9udC1zaXplPSIyMHB4IiBmaWxsPSIjMDAwMDAwIiB0ZXh0LWFuY2hvcj0ic3RhcnQiIHN0eWxlPSJ3aGl0ZS1zcGFjZTogcHJlOyIgZGlyZWN0aW9uPSJsdHIiPlN0b3JlcyBydWxlcywgPC90ZXh0Pjx0ZXh0IHg9IjAiIHk9IjQzIiBmb250LWZhbWlseT0iVmlyZ2lsLCBTZWdvZSBVSSBFbW9qaSIgZm9udC1zaXplPSIyMHB4IiBmaWxsPSIjMDAwMDAwIiB0ZXh0LWFuY2hvcj0ic3RhcnQiIHN0eWxlPSJ3aGl0ZS1zcGFjZTogcHJlOyIgZGlyZWN0aW9uPSJsdHIiPnN1YnNjcmlwdGlvbnMgJmFtcDs8L3RleHQ+PHRleHQgeD0iMCIgeT0iNjgiIGZvbnQtZmFtaWx5PSJWaXJnaWwsIFNlZ29lIFVJIEVtb2ppIiBmb250LXNpemU9IjIwcHgiIGZpbGw9IiMwMDAwMDAiIHRleHQtYW5jaG9yPSJzdGFydCIgc3R5bGU9IndoaXRlLXNwYWNlOiBwcmU7IiBkaXJlY3Rpb249Imx0ciI+YXV0aCBrZXlzPC90ZXh0PjwvZz48L3N2Zz4="}}]);