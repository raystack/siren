"use strict";(self.webpackChunksiren=self.webpackChunksiren||[]).push([[266],{3905:function(e,n,r){r.d(n,{Zo:function(){return c},kt:function(){return m}});var t=r(7294);function i(e,n,r){return n in e?Object.defineProperty(e,n,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[n]=r,e}function a(e,n){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var t=Object.getOwnPropertySymbols(e);n&&(t=t.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),r.push.apply(r,t)}return r}function o(e){for(var n=1;n<arguments.length;n++){var r=null!=arguments[n]?arguments[n]:{};n%2?a(Object(r),!0).forEach((function(n){i(e,n,r[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):a(Object(r)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(r,n))}))}return e}function s(e,n){if(null==e)return{};var r,t,i=function(e,n){if(null==e)return{};var r,t,i={},a=Object.keys(e);for(t=0;t<a.length;t++)r=a[t],n.indexOf(r)>=0||(i[r]=e[r]);return i}(e,n);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(t=0;t<a.length;t++)r=a[t],n.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(i[r]=e[r])}return i}var l=t.createContext({}),u=function(e){var n=t.useContext(l),r=n;return e&&(r="function"==typeof e?e(n):o(o({},n),e)),r},c=function(e){var n=u(e.components);return t.createElement(l.Provider,{value:n},e.children)},d={inlineCode:"code",wrapper:function(e){var n=e.children;return t.createElement(t.Fragment,{},n)}},p=t.forwardRef((function(e,n){var r=e.components,i=e.mdxType,a=e.originalType,l=e.parentName,c=s(e,["components","mdxType","originalType","parentName"]),p=u(r),m=i,g=p["".concat(l,".").concat(m)]||p[m]||d[m]||a;return r?t.createElement(g,o(o({ref:n},c),{},{components:r})):t.createElement(g,o({ref:n},c))}));function m(e,n){var r=arguments,i=n&&n.mdxType;if("string"==typeof e||i){var a=r.length,o=new Array(a);o[0]=p;var s={};for(var l in n)hasOwnProperty.call(n,l)&&(s[l]=n[l]);s.originalType=e,s.mdxType="string"==typeof e?e:i,o[1]=s;for(var u=2;u<a;u++)o[u]=r[u];return t.createElement.apply(null,o)}return t.createElement.apply(null,r)}p.displayName="MDXCreateElement"},4967:function(e,n,r){r.r(n),r.d(n,{frontMatter:function(){return s},contentTitle:function(){return l},metadata:function(){return u},toc:function(){return c},default:function(){return p}});var t=r(7462),i=r(3366),a=(r(7294),r(3905)),o=["components"],s={},l="Usage",u={unversionedId:"guides/overview",id:"guides/overview",isDocsHomePage:!1,title:"Usage",description:"The following topics will describe how to use Siren.",source:"@site/docs/guides/overview.md",sourceDirName:"guides",slug:"/guides/overview",permalink:"/siren/docs/guides/overview",editUrl:"https://github.com/odpf/siren/edit/master/docs/docs/guides/overview.md",tags:[],version:"current",lastUpdatedBy:"Abduh",lastUpdatedAt:1662977811,formattedLastUpdatedAt:"9/12/2022",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Introduction",permalink:"/siren/docs/introduction"},next:{title:"Providers",permalink:"/siren/docs/guides/providers"}},c=[{value:"CLI Interface",id:"cli-interface",children:[]},{value:"Managing providers and multi-tenancy",id:"managing-providers-and-multi-tenancy",children:[]},{value:"Managing Templates",id:"managing-templates",children:[]},{value:"Managing Rules",id:"managing-rules",children:[]},{value:"Managing bulk rules and templates",id:"managing-bulk-rules-and-templates",children:[]},{value:"Receivers",id:"receivers",children:[]},{value:"Subscriptions",id:"subscriptions",children:[]},{value:"Alert History Subscription",id:"alert-history-subscription",children:[]},{value:"Deployment",id:"deployment",children:[]},{value:"Monitoring",id:"monitoring",children:[]},{value:"Troubleshooting",id:"troubleshooting",children:[]}],d={toc:c};function p(e){var n=e.components,r=(0,i.Z)(e,o);return(0,a.kt)("wrapper",(0,t.Z)({},d,r,{components:n,mdxType:"MDXLayout"}),(0,a.kt)("h1",{id:"usage"},"Usage"),(0,a.kt)("p",null,"The following topics will describe how to use Siren."),(0,a.kt)("h2",{id:"cli-interface"},"CLI Interface"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-text"},"Siren provides alerting on metrics of your applications using Cortex metrics\nin a simple DIY configuration. With Siren, you can define templates(using go templates), and\ncreate/edit/enable/disable prometheus rules on demand.\n\nAvailable Commands:\n  alert       Manage alerts\n  completion  generate the autocompletion script for the specified shell\n  config      manage siren CLI configuration\n  help        Help about any command\n  migrate     Migrate database schema\n  namespace   Manage namespaces\n  provider    Manage providers\n  receiver    Manage receivers\n  rule        Manage rules\n  serve       Run server\n  template    Manage templates\n")),(0,a.kt)("h2",{id:"managing-providers-and-multi-tenancy"},"Managing providers and multi-tenancy"),(0,a.kt)("p",null,'Siren can be used define alerts and their routing configurations inside monitoring "providers". List of supported\nproviders:'),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},"CortexMetrics.")),(0,a.kt)("p",null,'Support for other providers is also planned, feel free to contribute. Siren also respects the multi-tenancy provided by\nvarious monitoring providers using "namespaces". Namespace simply represents a tenant inside your provider. Learn in\nmore detail ',(0,a.kt)("a",{parentName:"p",href:"/siren/docs/guides/providers"},"here"),"."),(0,a.kt)("h2",{id:"managing-templates"},"Managing Templates"),(0,a.kt)("p",null,"Siren templates are abstraction over Prometheus rules to reuse same rule body to create multiple rules. The rule body is\ntemplated using go templates. Learn in more detail ",(0,a.kt)("a",{parentName:"p",href:"/siren/docs/guides/templates"},"here"),"."),(0,a.kt)("h2",{id:"managing-rules"},"Managing Rules"),(0,a.kt)("p",null,"Siren rules are defined using a template by providing value for the variables defined inside that template. Learn in\nmore details ",(0,a.kt)("a",{parentName:"p",href:"/siren/docs/guides/rules"},"here")),(0,a.kt)("h2",{id:"managing-bulk-rules-and-templates"},"Managing bulk rules and templates"),(0,a.kt)("p",null,"For org wide use cases, where teams need to manage multiple templates and rules Siren CLI can be highly useful. Think\nGitOps but for alerting. Learn in More detail ",(0,a.kt)("a",{parentName:"p",href:"/siren/docs/guides/bulk_rules"},"here")),(0,a.kt)("h2",{id:"receivers"},"Receivers"),(0,a.kt)("p",null,"Receivers represent a notification medium, which can be used to define routing configuration in the monitoring\nproviders, to control the behaviour of how your alerts are notified. Few examples: Slack receiver, HTTP receiver,\nPagerduty receivers. You can use receivers to send notifications on demand as well as on certain matching conditions.\nLearn in more detail ",(0,a.kt)("a",{parentName:"p",href:"/siren/docs/guides/receivers"},"here"),"."),(0,a.kt)("h2",{id:"subscriptions"},"Subscriptions"),(0,a.kt)("p",null,"Siren can be used to configure various monitoring providers to route your alerts to proper channels based on your match\nconditions. You define your own set of selectors and subscribe to alerts matching these selectors in the notification\nmediums of your choice. Learn in more detail ",(0,a.kt)("a",{parentName:"p",href:"/siren/docs/guides/subscriptions"},"here"),"."),(0,a.kt)("h2",{id:"alert-history-subscription"},"Alert History Subscription"),(0,a.kt)("p",null,'Siren can configure Cortex Alertmanager to call Siren back, allowing storage of triggered alerts. This can be used for\nauditing and analytics purposes. Alert History is simply a "subscription" defined using an "HTTP receiver" on all\nalerts.'),(0,a.kt)("p",null,"Learn in more detail ",(0,a.kt)("a",{parentName:"p",href:"/siren/docs/guides/alert_history"},"here"),"."),(0,a.kt)("h2",{id:"deployment"},"Deployment"),(0,a.kt)("p",null,"Refer ",(0,a.kt)("a",{parentName:"p",href:"/siren/docs/guides/deployment"},"here")," to learn how to deploy Siren in production."),(0,a.kt)("h2",{id:"monitoring"},"Monitoring"),(0,a.kt)("p",null,"Refer ",(0,a.kt)("a",{parentName:"p",href:"/siren/docs/guides/monitoring"},"here")," to for more details on monitoring siren."),(0,a.kt)("h2",{id:"troubleshooting"},"Troubleshooting"),(0,a.kt)("p",null,"Troubleshooting ",(0,a.kt)("a",{parentName:"p",href:"/siren/docs/guides/troubleshooting"},"guide"),"."))}p.isMDXComponent=!0}}]);