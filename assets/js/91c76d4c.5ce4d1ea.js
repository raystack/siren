"use strict";(self.webpackChunksiren=self.webpackChunksiren||[]).push([[960],{3905:function(e,r,n){n.d(r,{Zo:function(){return c},kt:function(){return m}});var t=n(7294);function i(e,r,n){return r in e?Object.defineProperty(e,r,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[r]=n,e}function l(e,r){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var t=Object.getOwnPropertySymbols(e);r&&(t=t.filter((function(r){return Object.getOwnPropertyDescriptor(e,r).enumerable}))),n.push.apply(n,t)}return n}function a(e){for(var r=1;r<arguments.length;r++){var n=null!=arguments[r]?arguments[r]:{};r%2?l(Object(n),!0).forEach((function(r){i(e,r,n[r])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):l(Object(n)).forEach((function(r){Object.defineProperty(e,r,Object.getOwnPropertyDescriptor(n,r))}))}return e}function s(e,r){if(null==e)return{};var n,t,i=function(e,r){if(null==e)return{};var n,t,i={},l=Object.keys(e);for(t=0;t<l.length;t++)n=l[t],r.indexOf(n)>=0||(i[n]=e[n]);return i}(e,r);if(Object.getOwnPropertySymbols){var l=Object.getOwnPropertySymbols(e);for(t=0;t<l.length;t++)n=l[t],r.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(i[n]=e[n])}return i}var o=t.createContext({}),d=function(e){var r=t.useContext(o),n=r;return e&&(n="function"==typeof e?e(r):a(a({},r),e)),n},c=function(e){var r=d(e.components);return t.createElement(o.Provider,{value:r},e.children)},p={inlineCode:"code",wrapper:function(e){var r=e.children;return t.createElement(t.Fragment,{},r)}},u=t.forwardRef((function(e,r){var n=e.components,i=e.mdxType,l=e.originalType,o=e.parentName,c=s(e,["components","mdxType","originalType","parentName"]),u=d(n),m=i,f=u["".concat(o,".").concat(m)]||u[m]||p[m]||l;return n?t.createElement(f,a(a({ref:r},c),{},{components:n})):t.createElement(f,a({ref:r},c))}));function m(e,r){var n=arguments,i=r&&r.mdxType;if("string"==typeof e||i){var l=n.length,a=new Array(l);a[0]=u;var s={};for(var o in r)hasOwnProperty.call(r,o)&&(s[o]=r[o]);s.originalType=e,s.mdxType="string"==typeof e?e:i,a[1]=s;for(var d=2;d<l;d++)a[d]=n[d];return t.createElement.apply(null,a)}return t.createElement.apply(null,n)}u.displayName="MDXCreateElement"},2045:function(e,r,n){n.r(r),n.d(r,{assets:function(){return o},contentTitle:function(){return a},default:function(){return p},frontMatter:function(){return l},metadata:function(){return s},toc:function(){return d}});var t=n(3117),i=(n(7294),n(3905));const l={},a="CLI",s={unversionedId:"reference/cli",id:"reference/cli",title:"CLI",description:"siren alert",source:"@site/docs/reference/cli.md",sourceDirName:"reference",slug:"/reference/cli",permalink:"/siren/docs/reference/cli",draft:!1,editUrl:"https://github.com/odpf/siren/edit/master/docs/docs/reference/cli.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Receiver",permalink:"/siren/docs/reference/receiver"}},o={},d=[{value:"<code>siren alert</code>",id:"siren-alert",level:2},{value:"<code>siren alert list [flags]</code>",id:"siren-alert-list-flags",level:3},{value:"<code>siren completion [bash|zsh|fish|powershell]</code>",id:"siren-completion-bashzshfishpowershell",level:2},{value:"<code>siren config &lt;command&gt;</code>",id:"siren-config-command",level:2},{value:"<code>siren config init</code>",id:"siren-config-init",level:3},{value:"<code>siren config list</code>",id:"siren-config-list",level:3},{value:"<code>siren environment</code>",id:"siren-environment",level:2},{value:"<code>siren job &lt;command&gt;</code>",id:"siren-job-command",level:2},{value:"<code>siren job run</code>",id:"siren-job-run",level:3},{value:"<code>siren job run cleanup_queue [flags]</code>",id:"siren-job-run-cleanup_queue-flags",level:4},{value:"<code>siren namespace</code>",id:"siren-namespace",level:2},{value:"<code>siren namespace create [flags]</code>",id:"siren-namespace-create-flags",level:3},{value:"<code>siren namespace delete</code>",id:"siren-namespace-delete",level:3},{value:"<code>siren namespace edit [flags]</code>",id:"siren-namespace-edit-flags",level:3},{value:"<code>siren namespace list</code>",id:"siren-namespace-list",level:3},{value:"<code>siren namespace view [flags]</code>",id:"siren-namespace-view-flags",level:3},{value:"<code>siren provider</code>",id:"siren-provider",level:2},{value:"<code>siren provider create [flags]</code>",id:"siren-provider-create-flags",level:3},{value:"<code>siren provider delete</code>",id:"siren-provider-delete",level:3},{value:"<code>siren provider edit [flags]</code>",id:"siren-provider-edit-flags",level:3},{value:"<code>siren provider list</code>",id:"siren-provider-list",level:3},{value:"<code>siren provider view [flags]</code>",id:"siren-provider-view-flags",level:3},{value:"<code>siren receiver</code>",id:"siren-receiver",level:2},{value:"<code>siren receiver create [flags]</code>",id:"siren-receiver-create-flags",level:3},{value:"<code>siren receiver delete</code>",id:"siren-receiver-delete",level:3},{value:"<code>siren receiver edit [flags]</code>",id:"siren-receiver-edit-flags",level:3},{value:"<code>siren receiver list</code>",id:"siren-receiver-list",level:3},{value:"<code>siren receiver send [flags]</code>",id:"siren-receiver-send-flags",level:3},{value:"<code>siren receiver view [flags]</code>",id:"siren-receiver-view-flags",level:3},{value:"<code>siren rule</code>",id:"siren-rule",level:2},{value:"<code>siren rule edit [flags]</code>",id:"siren-rule-edit-flags",level:3},{value:"<code>siren rule list [flags]</code>",id:"siren-rule-list-flags",level:3},{value:"<code>siren rule upload</code>",id:"siren-rule-upload",level:3},{value:"<code>siren server &lt;command&gt;</code>",id:"siren-server-command",level:2},{value:"<code>siren server init [flags]</code>",id:"siren-server-init-flags",level:3},{value:"<code>siren server migrate [flags]</code>",id:"siren-server-migrate-flags",level:3},{value:"<code>siren server start [flags]</code>",id:"siren-server-start-flags",level:3},{value:"<code>siren template</code>",id:"siren-template",level:2},{value:"<code>siren template delete</code>",id:"siren-template-delete",level:3},{value:"<code>siren template list [flags]</code>",id:"siren-template-list-flags",level:3},{value:"<code>siren template render [flags]</code>",id:"siren-template-render-flags",level:3},{value:"<code>siren template upload</code>",id:"siren-template-upload",level:3},{value:"<code>siren template upsert [flags]</code>",id:"siren-template-upsert-flags",level:3},{value:"<code>siren template view [flags]</code>",id:"siren-template-view-flags",level:3},{value:"<code>siren worker &lt;command&gt; &lt;worker_command&gt;</code>",id:"siren-worker-command-worker_command",level:2},{value:"<code>siren worker start &lt;command&gt;</code>",id:"siren-worker-start-command",level:3},{value:"<code>siren worker start notification_dlq_handler [flags]</code>",id:"siren-worker-start-notification_dlq_handler-flags",level:4},{value:"<code>siren worker start notification_handler [flags]</code>",id:"siren-worker-start-notification_handler-flags",level:4}],c={toc:d};function p(e){let{components:r,...n}=e;return(0,i.kt)("wrapper",(0,t.Z)({},c,n,{components:r,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"cli"},"CLI"),(0,i.kt)("h2",{id:"siren-alert"},(0,i.kt)("inlineCode",{parentName:"h2"},"siren alert")),(0,i.kt)("p",null,"Manage alerts"),(0,i.kt)("h3",{id:"siren-alert-list-flags"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren alert list [flags]")),(0,i.kt)("p",null,"List alerts"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},"--end-time uint          end time\n--provider-id uint       provider id\n--provider-type string   provider type\n--resource-name string   resource name\n--start-time uint        start time\n")),(0,i.kt)("h2",{id:"siren-completion-bashzshfishpowershell"},(0,i.kt)("inlineCode",{parentName:"h2"},"siren completion [bash|zsh|fish|powershell]")),(0,i.kt)("p",null,"Generate shell completion scripts"),(0,i.kt)("h2",{id:"siren-config-command"},(0,i.kt)("inlineCode",{parentName:"h2"},"siren config <command>")),(0,i.kt)("p",null,"Manage siren CLI configuration"),(0,i.kt)("h3",{id:"siren-config-init"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren config init")),(0,i.kt)("p",null,"Initialize CLI configuration"),(0,i.kt)("h3",{id:"siren-config-list"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren config list")),(0,i.kt)("p",null,"List client configuration settings"),(0,i.kt)("h2",{id:"siren-environment"},(0,i.kt)("inlineCode",{parentName:"h2"},"siren environment")),(0,i.kt)("p",null,"List of supported environment variables"),(0,i.kt)("h2",{id:"siren-job-command"},(0,i.kt)("inlineCode",{parentName:"h2"},"siren job <command>")),(0,i.kt)("p",null,"Manage siren jobs"),(0,i.kt)("h3",{id:"siren-job-run"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren job run")),(0,i.kt)("p",null,"Trigger a job"),(0,i.kt)("h4",{id:"siren-job-run-cleanup_queue-flags"},(0,i.kt)("inlineCode",{parentName:"h4"},"siren job run cleanup_queue [flags]")),(0,i.kt)("p",null,"Cleanup stale messages in queue"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},'-c, --config string   Config file path (default "config.yaml")\n')),(0,i.kt)("h2",{id:"siren-namespace"},(0,i.kt)("inlineCode",{parentName:"h2"},"siren namespace")),(0,i.kt)("p",null,"Manage namespaces"),(0,i.kt)("h3",{id:"siren-namespace-create-flags"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren namespace create [flags]")),(0,i.kt)("p",null,"Create a new namespace"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},"-f, --file string   path to the namespace config\n")),(0,i.kt)("h3",{id:"siren-namespace-delete"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren namespace delete")),(0,i.kt)("p",null,"Delete a namespace details"),(0,i.kt)("h3",{id:"siren-namespace-edit-flags"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren namespace edit [flags]")),(0,i.kt)("p",null,"Edit a namespace"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},"-f, --file string   Path to the namespace config\n    --id uint       namespace id\n")),(0,i.kt)("h3",{id:"siren-namespace-list"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren namespace list")),(0,i.kt)("p",null,"List namespaces"),(0,i.kt)("h3",{id:"siren-namespace-view-flags"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren namespace view [flags]")),(0,i.kt)("p",null,"View a namespace details"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},'--format string   Print output with the selected format (default "yaml")\n')),(0,i.kt)("h2",{id:"siren-provider"},(0,i.kt)("inlineCode",{parentName:"h2"},"siren provider")),(0,i.kt)("p",null,"Manage providers"),(0,i.kt)("h3",{id:"siren-provider-create-flags"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren provider create [flags]")),(0,i.kt)("p",null,"Create a new provider"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},"-f, --file string   path to the provider config\n")),(0,i.kt)("h3",{id:"siren-provider-delete"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren provider delete")),(0,i.kt)("p",null,"Delete a provider details"),(0,i.kt)("h3",{id:"siren-provider-edit-flags"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren provider edit [flags]")),(0,i.kt)("p",null,"Edit a provider"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},"-f, --file string   Path to the provider config\n    --id uint       provider id\n")),(0,i.kt)("h3",{id:"siren-provider-list"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren provider list")),(0,i.kt)("p",null,"List providers"),(0,i.kt)("h3",{id:"siren-provider-view-flags"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren provider view [flags]")),(0,i.kt)("p",null,"View a provider details"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},'--format string   Print output with the selected format (default "yaml")\n')),(0,i.kt)("h2",{id:"siren-receiver"},(0,i.kt)("inlineCode",{parentName:"h2"},"siren receiver")),(0,i.kt)("p",null,"Manage receivers"),(0,i.kt)("h3",{id:"siren-receiver-create-flags"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren receiver create [flags]")),(0,i.kt)("p",null,"Create a new receiver"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},"-f, --file string   path to the receiver config\n")),(0,i.kt)("h3",{id:"siren-receiver-delete"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren receiver delete")),(0,i.kt)("p",null,"Delete a receiver details"),(0,i.kt)("h3",{id:"siren-receiver-edit-flags"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren receiver edit [flags]")),(0,i.kt)("p",null,"Edit a receiver"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},"-f, --file string   Path to the receiver config\n    --id uint       receiver id\n")),(0,i.kt)("h3",{id:"siren-receiver-list"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren receiver list")),(0,i.kt)("p",null,"List receivers"),(0,i.kt)("h3",{id:"siren-receiver-send-flags"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren receiver send [flags]")),(0,i.kt)("p",null,"Send a receiver notification"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},"-f, --file string   Path to the receiver notification message\n    --id uint       receiver id\n")),(0,i.kt)("h3",{id:"siren-receiver-view-flags"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren receiver view [flags]")),(0,i.kt)("p",null,"View a receiver details"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},'--format string   Print output with the selected format (default "yaml")\n')),(0,i.kt)("h2",{id:"siren-rule"},(0,i.kt)("inlineCode",{parentName:"h2"},"siren rule")),(0,i.kt)("p",null,"Manage rules"),(0,i.kt)("h3",{id:"siren-rule-edit-flags"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren rule edit [flags]")),(0,i.kt)("p",null,"Edit a rule"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},"-f, --file string   Path to the rule config\n    --id uint       rule id\n")),(0,i.kt)("h3",{id:"siren-rule-list-flags"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren rule list [flags]")),(0,i.kt)("p",null,"List rules"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},"--group-name string         rule group name\n--name string               rule name\n--namespace string          rule namespace\n--provider-namespace uint   rule provider namespace id\n--template string           rule template\n")),(0,i.kt)("h3",{id:"siren-rule-upload"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren rule upload")),(0,i.kt)("p",null,"Upload Rules YAML file"),(0,i.kt)("h2",{id:"siren-server-command"},(0,i.kt)("inlineCode",{parentName:"h2"},"siren server <command>")),(0,i.kt)("p",null,"Run siren server"),(0,i.kt)("h3",{id:"siren-server-init-flags"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren server init [flags]")),(0,i.kt)("p",null,"Initialize server"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},'-o, --output string   Output config file path (default "./config.yaml")\n')),(0,i.kt)("h3",{id:"siren-server-migrate-flags"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren server migrate [flags]")),(0,i.kt)("p",null,"Run DB Schema Migrations"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},'-c, --config string   Config file path (default "./config.yaml")\n')),(0,i.kt)("h3",{id:"siren-server-start-flags"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren server start [flags]")),(0,i.kt)("p",null,"Start server on default port 8080"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},'-c, --config string   Config file path (default "config.yaml")\n')),(0,i.kt)("h2",{id:"siren-template"},(0,i.kt)("inlineCode",{parentName:"h2"},"siren template")),(0,i.kt)("p",null,"Manage templates"),(0,i.kt)("h3",{id:"siren-template-delete"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren template delete")),(0,i.kt)("p",null,"Delete a template details"),(0,i.kt)("h3",{id:"siren-template-list-flags"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren template list [flags]")),(0,i.kt)("p",null,"List templates"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},"--tag string   template tag name\n")),(0,i.kt)("h3",{id:"siren-template-render-flags"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren template render [flags]")),(0,i.kt)("p",null,"Render a template details"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},'-f, --file string     path to the template config\n    --format string   Print output with the selected format (default "yaml")\n    --name string     template name\n')),(0,i.kt)("h3",{id:"siren-template-upload"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren template upload")),(0,i.kt)("p",null,"Upload Templates YAML file"),(0,i.kt)("h3",{id:"siren-template-upsert-flags"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren template upsert [flags]")),(0,i.kt)("p",null,"Create or edit a new template"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},"-f, --file string   path to the template config\n")),(0,i.kt)("h3",{id:"siren-template-view-flags"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren template view [flags]")),(0,i.kt)("p",null,"View a template details"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},'--format string   Print output with the selected format (default "yaml")\n')),(0,i.kt)("h2",{id:"siren-worker-command-worker_command"},(0,i.kt)("inlineCode",{parentName:"h2"},"siren worker <command> <worker_command>")),(0,i.kt)("p",null,"Start or manage Siren's workers"),(0,i.kt)("h3",{id:"siren-worker-start-command"},(0,i.kt)("inlineCode",{parentName:"h3"},"siren worker start <command>")),(0,i.kt)("p",null,"Start a siren worker"),(0,i.kt)("h4",{id:"siren-worker-start-notification_dlq_handler-flags"},(0,i.kt)("inlineCode",{parentName:"h4"},"siren worker start notification_dlq_handler [flags]")),(0,i.kt)("p",null,"A notification dlq handler"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},'-c, --config string   Config file path (default "config.yaml")\n')),(0,i.kt)("h4",{id:"siren-worker-start-notification_handler-flags"},(0,i.kt)("inlineCode",{parentName:"h4"},"siren worker start notification_handler [flags]")),(0,i.kt)("p",null,"A notification handler"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},'-c, --config string   Config file path (default "config.yaml")\n')))}p.isMDXComponent=!0}}]);