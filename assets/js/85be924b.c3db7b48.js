"use strict";(self.webpackChunksiren=self.webpackChunksiren||[]).push([[861],{3905:function(e,t,n){n.d(t,{Zo:function(){return p},kt:function(){return d}});var i=n(7294);function r(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);t&&(i=i.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,i)}return n}function a(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){r(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function s(e,t){if(null==e)return{};var n,i,r=function(e,t){if(null==e)return{};var n,i,r={},o=Object.keys(e);for(i=0;i<o.length;i++)n=o[i],t.indexOf(n)>=0||(r[n]=e[n]);return r}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(i=0;i<o.length;i++)n=o[i],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(r[n]=e[n])}return r}var l=i.createContext({}),c=function(e){var t=i.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):a(a({},t),e)),n},p=function(e){var t=c(e.components);return i.createElement(l.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return i.createElement(i.Fragment,{},t)}},f=i.forwardRef((function(e,t){var n=e.components,r=e.mdxType,o=e.originalType,l=e.parentName,p=s(e,["components","mdxType","originalType","parentName"]),f=c(n),d=r,g=f["".concat(l,".").concat(d)]||f[d]||u[d]||o;return n?i.createElement(g,a(a({ref:t},p),{},{components:n})):i.createElement(g,a({ref:t},p))}));function d(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var o=n.length,a=new Array(o);a[0]=f;var s={};for(var l in t)hasOwnProperty.call(t,l)&&(s[l]=t[l]);s.originalType=e,s.mdxType="string"==typeof e?e:r,a[1]=s;for(var c=2;c<o;c++)a[c]=n[c];return i.createElement.apply(null,a)}return i.createElement.apply(null,n)}f.displayName="MDXCreateElement"},7825:function(e,t,n){n.r(t),n.d(t,{assets:function(){return l},contentTitle:function(){return a},default:function(){return u},frontMatter:function(){return o},metadata:function(){return s},toc:function(){return c}});var i=n(3117),r=(n(7294),n(3905));const o={},a="Plugin",s={unversionedId:"concepts/plugin",id:"concepts/plugin",title:"Plugin",description:"Siren decouples various provider, receiver, and queue as a plugin. The purpose is to ease the extension of new plugin. We welcome all contributions to add new plugin.",source:"@site/docs/concepts/plugin.md",sourceDirName:"concepts",slug:"/concepts/plugin",permalink:"/siren/docs/concepts/plugin",draft:!1,editUrl:"https://github.com/odpf/siren/edit/master/docs/docs/concepts/plugin.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Overview",permalink:"/siren/docs/concepts/overview"},next:{title:"Schema Design",permalink:"/siren/docs/concepts/schema"}},l={},c=[{value:"Provider",id:"provider",level:2},{value:"Receiver",id:"receiver",level:2},{value:"Configurations",id:"configurations",level:3},{value:"Interface",id:"interface",level:3},{value:"ConfigResolver",id:"configresolver",level:4},{value:"Notifier",id:"notifier",level:4},{value:"Base Plugin",id:"base-plugin",level:3},{value:"Queue",id:"queue",level:2}],p={toc:c};function u(e){let{components:t,...n}=e;return(0,r.kt)("wrapper",(0,i.Z)({},p,n,{components:t,mdxType:"MDXLayout"}),(0,r.kt)("h1",{id:"plugin"},"Plugin"),(0,r.kt)("p",null,"Siren decouples various ",(0,r.kt)("inlineCode",{parentName:"p"},"provider"),", ",(0,r.kt)("inlineCode",{parentName:"p"},"receiver"),", and ",(0,r.kt)("inlineCode",{parentName:"p"},"queue")," as a plugin. The purpose is to ease the extension of new plugin. We welcome all contributions to add new plugin."),(0,r.kt)("h2",{id:"provider"},"Provider"),(0,r.kt)("p",null,"Provider responsibility is to accept incoming rules configuration from Siren and send alerts to the designated Siren Hook API. Provider plugin needs to fulfill some interfaces. More detail about interfaces can be found in ",(0,r.kt)("a",{parentName:"p",href:"/siren/docs/contribute/provider"},"contribution")," page. Supported providers are:"),(0,r.kt)("ul",null,(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("a",{parentName:"li",href:"https://cortexmetrics.io/"},"Cortexmetrics"))),(0,r.kt)("h2",{id:"receiver"},"Receiver"),(0,r.kt)("p",null,"Receiver defines where the notification Siren sends to. Receiver plugin needs to fulfill some interfaces. More detail about interfaces can be found in ",(0,r.kt)("a",{parentName:"p",href:"/siren/docs/contribute/receiver"},"contribution")," page. Supported providers are:"),(0,r.kt)("ul",null,(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("a",{parentName:"li",href:"https://api.slack.com/methods/chat.postMessage"},"Slack")),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("a",{parentName:"li",href:"https://developer.pagerduty.com/docs/ZG9jOjExMDI5NTc3-events-api-v1"},"PagerDuty Events API v1")),(0,r.kt)("li",{parentName:"ul"},"HTTP"),(0,r.kt)("li",{parentName:"ul"},"File")),(0,r.kt)("p",null,"Receiver plugin is being used by two different services: receiver service and notification service. Receiver service handles the way the receiver is being stored, updated, fetched, and removed. Notification service uses receiver plugin to send notification. Each service has its own interface that needs to be implemented."),(0,r.kt)("h3",{id:"configurations"},"Configurations"),(0,r.kt)("p",null,"Siren receiver plugins have several configs: ",(0,r.kt)("inlineCode",{parentName:"p"},"ReceiverConfig"),", ",(0,r.kt)("inlineCode",{parentName:"p"},"SubscriptionConfig")," (if needed), and ",(0,r.kt)("inlineCode",{parentName:"p"},"NotificationConfig"),"."),(0,r.kt)("ul",null,(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("strong",{parentName:"li"},"ReceiverConfig")," is a config that will be part of ",(0,r.kt)("inlineCode",{parentName:"li"},"receiver.Receiver")," struct and will be stored inside the DB's receivers table."),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("strong",{parentName:"li"},"SubscriptionConfig")," is optional. Subscription config is defined and used if the receivers inside subscription requires another additional configs rather than ",(0,r.kt)("inlineCode",{parentName:"li"},"ReceiverConfig"),". For example, Slack stores encrypted ",(0,r.kt)("inlineCode",{parentName:"li"},"token")," when storing receiver information inside the DB but has another config ",(0,r.kt)("inlineCode",{parentName:"li"},"channel_name")," on subscription level."),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("strong",{parentName:"li"},"NotificationConfig")," embeds ",(0,r.kt)("inlineCode",{parentName:"li"},"ReceiverConfig")," and ",(0,r.kt)("inlineCode",{parentName:"li"},"SubscriptionConfig")," (if needed)."),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("strong",{parentName:"li"},"AppConfig")," is a config of receiver plugins that is being loaded when the Siren app is started. ",(0,r.kt)("inlineCode",{parentName:"li"},"AppConfig")," can be set up via environment variable or config file. Usually this is a generic config of a specific receiver regardless where the notification is being sent to (e.g. http config, receiver host, etc...). If your plugin requires ",(0,r.kt)("inlineCode",{parentName:"li"},"AppConfig"),", you can set the config inside ",(0,r.kt)("inlineCode",{parentName:"li"},"plugins/receivers/config.go"),".")),(0,r.kt)("p",null,"In Siren receiver plugins, all configs will be transform back and forth from ",(0,r.kt)("inlineCode",{parentName:"p"},"map[string]interface{}")," to struct using ",(0,r.kt)("a",{parentName:"p",href:"https://github.com/mitchellh/mapstructure"},"mitchellh/mapstructure"),". You might also need to add more functions to validate and transform configs to ",(0,r.kt)("inlineCode",{parentName:"p"},"map[string]interface{}"),"."),(0,r.kt)("h3",{id:"interface"},"Interface"),(0,r.kt)("h4",{id:"configresolver"},"ConfigResolver"),(0,r.kt)("p",null,"ConfigResolver is being used by receiver service to manage receivers. It is an interface for the receiver to resolve all configs and functions."),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-go"},"type ConfigResolver interface {\n    // TODO might be removed\n    BuildData(ctx context.Context, configs map[string]interface{}) (map[string]interface{}, error)\n    // TODO might be removed\n    BuildNotificationConfig(subscriptionConfigMap map[string]interface{}, receiverConfigMap map[string]interface{}) (map[string]interface{}, error)\n    PreHookTransformConfigs(ctx context.Context, configs map[string]interface{}) (map[string]interface{}, error)\n    PostHookTransformConfigs(ctx context.Context, configs map[string]interface{}) (map[string]interface{}, error)\n}\n")),(0,r.kt)("ul",null,(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("strong",{parentName:"li"},"BuildData")," is being used in GetReceiver where ",(0,r.kt)("inlineCode",{parentName:"li"},"data")," field in Receiver is being populated. This might not relevant anymore for our current use case and might be deprecated later."),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("strong",{parentName:"li"},"BuildNotificationConfig")," is being used for subscription. This might not relevant anymore for our current use case and might be deprecated later."),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("strong",{parentName:"li"},"PreHookTransformConfigs")," is being used to transform configs (e.g. encryption) before the config is being stored in the DB."),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("strong",{parentName:"li"},"PostHookTransformConfigs")," is being used to transform configs (e.g. decryption) after the config is being fetched from the DB.")),(0,r.kt)("h4",{id:"notifier"},"Notifier"),(0,r.kt)("p",null,"Notifier interface is being used by notification service and consists of all functionalities to publish notifications."),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-go"},"type Notifier interface {\n    PreHookTransformConfigs(ctx context.Context, notificationConfigMap map[string]interface{}) (map[string]interface{}, error)\n    PostHookTransformConfigs(ctx context.Context, notificationConfigMap map[string]interface{}) (map[string]interface{}, error)\n    DefaultTemplateOfProvider(templateName string) string\n    Publish(ctx context.Context, message Message) (bool, error)\n}\n")),(0,r.kt)("ul",null,(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("strong",{parentName:"li"},"PreHookTransformConfigs")," is being used to transform configs (e.g. encryption) before the config is being enqueued."),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("strong",{parentName:"li"},"PostHookTransformConfigs")," is being used to transform configs (e.g. decryption) after the config is being dequeued."),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("strong",{parentName:"li"},"DefaultTemplateOfProvider")," assigns default provider template for alert notifications of as specific provider. Each provider might send alerts with different format, the template needs to build notification specific message out of the alerts for each provider. Each provider has to have a reserved template name (e.g. ",(0,r.kt)("inlineCode",{parentName:"li"},"template.ReservedName_xxx"),") and all alerts coming from the provider needs to use the template with the reserved name."),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("strong",{parentName:"li"},"Publish")," handles how message is being sent. The first return argument is ",(0,r.kt)("inlineCode",{parentName:"li"},"retryable")," boolean to indicate whether an error is a ",(0,r.kt)("inlineCode",{parentName:"li"},"retryable")," error or not. If it is API call, usually response status code 429 or 5xx is retriable. You can use ",(0,r.kt)("inlineCode",{parentName:"li"},"pkg/retrier")," to retry the call.")),(0,r.kt)("h3",{id:"base-plugin"},"Base Plugin"),(0,r.kt)("p",null,"Siren provide base plugin in ",(0,r.kt)("inlineCode",{parentName:"p"},"plugins/receivers/base")," which can be embedded in all plugins service struct. By doing so, you just need to implement all interfaces' method that you only need. The unimplemented methods one will already be handled by the ",(0,r.kt)("inlineCode",{parentName:"p"},"base")," plugin."),(0,r.kt)("h2",{id:"queue"},"Queue"),(0,r.kt)("p",null,"Queue is used as a buffer for the outbound notifications. Siren has a pluggable queue where user could choose which Queue to use in the ",(0,r.kt)("a",{parentName:"p",href:"/siren/docs/reference/server_configuration"},"config"),". Supported Queues are:"),(0,r.kt)("ul",null,(0,r.kt)("li",{parentName:"ul"},"In-Memory"),(0,r.kt)("li",{parentName:"ul"},"PostgreSQl")))}u.isMDXComponent=!0}}]);