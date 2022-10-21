"use strict";(self.webpackChunksiren=self.webpackChunksiren||[]).push([[463],{7313:function(e,r,t){t.r(r),t.d(r,{frontMatter:function(){return a},contentTitle:function(){return u},metadata:function(){return l},toc:function(){return d},default:function(){return h}});var n=t(7462),s=t(3366),i=(t(7294),t(3905)),o=["components"],a={},u="Workers",l={unversionedId:"guides/worker",id:"guides/worker",isDocsHomePage:!1,title:"Workers",description:"Siren has a notification features that utilizes queue to publish notification messages. The architecture requires a detached worker running asynchronously and polling queue periodically to dequeue notification messages and publish them. By default, Siren server run this asynchronous worker inside it. However it is also possible to run the worker as a different process. Currently there are two possible workers to run",source:"@site/docs/guides/worker.md",sourceDirName:"guides",slug:"/guides/worker",permalink:"/siren/docs/guides/worker",editUrl:"https://github.com/odpf/siren/edit/master/docs/docs/guides/worker.md",tags:[],version:"current",lastUpdatedBy:"Ravi Suhag",lastUpdatedAt:1666359878,formattedLastUpdatedAt:"10/21/2022",frontMatter:{}},d=[{value:"Running Workers as a Different Process",id:"running-workers-as-a-different-process",children:[]}],c={toc:d};function h(e){var r=e.components,t=(0,s.Z)(e,o);return(0,i.kt)("wrapper",(0,n.Z)({},c,t,{components:r,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"workers"},"Workers"),(0,i.kt)("p",null,"Siren has a notification features that utilizes queue to publish notification messages. The architecture requires a detached worker running asynchronously and polling queue periodically to dequeue notification messages and publish them. By default, Siren server run this asynchronous worker inside it. However it is also possible to run the worker as a different process. Currently there are two possible workers to run"),(0,i.kt)("ol",null,(0,i.kt)("li",{parentName:"ol"},(0,i.kt)("strong",{parentName:"li"},"Notification message handler:")," this worker periodically poll and dequeue messages from queue, process the messages, and then publish notification messages to the specified receivers. If there is a failure, this handler enqueues the failed messages to the dlq."),(0,i.kt)("li",{parentName:"ol"},(0,i.kt)("strong",{parentName:"li"},"Notification dlq handler:")," this worker periodically poll and dequeue messages from dead-letter-queue, process the messages, and then publish notification messages to the specified receivers. If there is a failure, this handler enqueues the failed messages back to the dlq.")),(0,i.kt)("h2",{id:"running-workers-as-a-different-process"},"Running Workers as a Different Process"),(0,i.kt)("p",null,"Siren has a command to start workers. Workers use the same config like server does."),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-bash"},"Command to start a siren worker.\n\nUsage\n  siren worker start <command> [flags]\n\nCore commands\n  notification_dlq_handler    A notification dlq handler\n  notification_handler        A notification handler\n\nInherited flags\n  --help   Show help for command\n\nExamples\n  $ siren worker start notification_handler\n  $ siren server start notification_handler -c ./config.yaml\n")),(0,i.kt)("p",null,"Starting up a worker could be done by executing."),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-bash"},"$ siren worker start notification_handler\n")))}h.isMDXComponent=!0}}]);