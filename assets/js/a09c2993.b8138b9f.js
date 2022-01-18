"use strict";(self.webpackChunksiren=self.webpackChunksiren||[]).push([[128],{3905:function(e,t,r){r.d(t,{Zo:function(){return u},kt:function(){return m}});var n=r(7294);function i(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function a(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function o(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?a(Object(r),!0).forEach((function(t){i(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):a(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function s(e,t){if(null==e)return{};var r,n,i=function(e,t){if(null==e)return{};var r,n,i={},a=Object.keys(e);for(n=0;n<a.length;n++)r=a[n],t.indexOf(r)>=0||(i[r]=e[r]);return i}(e,t);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(n=0;n<a.length;n++)r=a[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(i[r]=e[r])}return i}var c=n.createContext({}),l=function(e){var t=n.useContext(c),r=t;return e&&(r="function"==typeof e?e(t):o(o({},t),e)),r},u=function(e){var t=l(e.components);return n.createElement(c.Provider,{value:t},e.children)},p={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},d=n.forwardRef((function(e,t){var r=e.components,i=e.mdxType,a=e.originalType,c=e.parentName,u=s(e,["components","mdxType","originalType","parentName"]),d=l(r),m=i,f=d["".concat(c,".").concat(m)]||d[m]||p[m]||a;return r?n.createElement(f,o(o({ref:t},u),{},{components:r})):n.createElement(f,o({ref:t},u))}));function m(e,t){var r=arguments,i=t&&t.mdxType;if("string"==typeof e||i){var a=r.length,o=new Array(a);o[0]=d;var s={};for(var c in t)hasOwnProperty.call(t,c)&&(s[c]=t[c]);s.originalType=e,s.mdxType="string"==typeof e?e:i,o[1]=s;for(var l=2;l<a;l++)o[l]=r[l];return n.createElement.apply(null,o)}return n.createElement.apply(null,r)}d.displayName="MDXCreateElement"},8495:function(e,t,r){r.r(t),r.d(t,{frontMatter:function(){return s},contentTitle:function(){return c},metadata:function(){return l},toc:function(){return u},default:function(){return d}});var n=r(7462),i=r(3366),a=(r(7294),r(3905)),o=["components"],s={},c="Introduction",l={unversionedId:"introduction",id:"introduction",isDocsHomePage:!1,title:"Introduction",description:"Siren provides alerting on metrics of your applications using Cortex metrics in a simple",source:"@site/docs/introduction.md",sourceDirName:".",slug:"/introduction",permalink:"/siren/docs/introduction",editUrl:"https://github.com/odpf/siren/edit/master/docs/docs/introduction.md",tags:[],version:"current",lastUpdatedBy:"Ravi Suhag",lastUpdatedAt:1642494349,formattedLastUpdatedAt:"1/18/2022",frontMatter:{},sidebar:"docsSidebar",next:{title:"Usage",permalink:"/siren/docs/guides/overview"}},u=[{value:"Key Features",id:"key-features",children:[]},{value:"Usage",id:"usage",children:[]}],p={toc:u};function d(e){var t=e.components,s=(0,i.Z)(e,o);return(0,a.kt)("wrapper",(0,n.Z)({},p,s,{components:t,mdxType:"MDXLayout"}),(0,a.kt)("h1",{id:"introduction"},"Introduction"),(0,a.kt)("p",null,"Siren provides alerting on metrics of your applications using ",(0,a.kt)("a",{parentName:"p",href:"https://cortexmetrics.io/"},"Cortex metrics")," in a simple\nDIY configuration. With Siren, you can define templates(using go templates standard), and create/edit/enable/disable\nprometheus rules on demand. It also gives flexibility to manage bulk of rules via YAML files. Siren can be integrated\nwith any client such as CI/CD pipelines, Self-Serve UI, microservices etc."),(0,a.kt)("p",null,(0,a.kt)("img",{alt:"Siren Alert History",src:r(4454).Z})),(0,a.kt)("h2",{id:"key-features"},"Key Features"),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("strong",{parentName:"li"},"Rule Templates:")," Siren provides a way to define templates over prometheus Rule, which can be reused to create\nmultiple instances of same rule with configurable thresholds."),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("strong",{parentName:"li"},"Multi-tenancy:")," Rules created with Siren are by default multi-tenancy aware."),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("strong",{parentName:"li"},"DIY Interface:")," Siren can be used to easily create/edit prometheus rules. It also provides soft delete(disable)\nso that you can preserve thresholds in case you need to reuse the same alert."),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("strong",{parentName:"li"},"Managing bulk rules:")," Siren enables users to manage bulk alerts using YAML files in specified format using simple\nCLI."),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("strong",{parentName:"li"},"Receivers:")," Siren can be used to send out notifications via several mediums(e.g. slack, pagerduty, email etc)."),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("strong",{parentName:"li"},"Subscriptions")," Siren can be used to subscribe to alerts (with desired matching conditions) via the channel of your\nchoice."),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("strong",{parentName:"li"},"Alert History:")," Siren can store alerts triggered via Cortex Alertmanager, which can be used for audit purposes.")),(0,a.kt)("h2",{id:"usage"},"Usage"),(0,a.kt)("p",null,"Explore the following resources to get started with Siren:"),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("a",{parentName:"li",href:"/siren/docs/guides/overview"},"Guides")," provides guidance on usage."),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("a",{parentName:"li",href:"/siren/docs/concepts/overview"},"Concepts")," describes all important Siren concepts including system architecture."),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("a",{parentName:"li",href:"/siren/docs/reference/configuration"},"Reference")," contains the details about configurations and other aspects of Siren."),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("a",{parentName:"li",href:"/siren/docs/contribute/contribution"},"Contribute")," contains resources for anyone who wants to contribute to Siren.")))}d.isMDXComponent=!0},4454:function(e,t,r){t.Z=r.p+"assets/images/overview-640dcd08ea55323369bae78b6055feec.svg"}}]);