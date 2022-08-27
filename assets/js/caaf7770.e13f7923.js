"use strict";(self.webpackChunksiren=self.webpackChunksiren||[]).push([[119],{3905:function(e,t,n){n.d(t,{Zo:function(){return p},kt:function(){return m}});var r=n(7294);function o(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function a(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){o(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function u(e,t){if(null==e)return{};var n,r,o=function(e,t){if(null==e)return{};var n,r,o={},i=Object.keys(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||(o[n]=e[n]);return o}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(o[n]=e[n])}return o}var l=r.createContext({}),c=function(e){var t=r.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):a(a({},t),e)),n},p=function(e){var t=c(e.components);return r.createElement(l.Provider,{value:t},e.children)},s={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},d=r.forwardRef((function(e,t){var n=e.components,o=e.mdxType,i=e.originalType,l=e.parentName,p=u(e,["components","mdxType","originalType","parentName"]),d=c(n),m=o,f=d["".concat(l,".").concat(m)]||d[m]||s[m]||i;return n?r.createElement(f,a(a({ref:t},p),{},{components:n})):r.createElement(f,a({ref:t},p))}));function m(e,t){var n=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var i=n.length,a=new Array(i);a[0]=d;var u={};for(var l in t)hasOwnProperty.call(t,l)&&(u[l]=t[l]);u.originalType=e,u.mdxType="string"==typeof e?e:o,a[1]=u;for(var c=2;c<i;c++)a[c]=n[c];return r.createElement.apply(null,a)}return r.createElement.apply(null,n)}d.displayName="MDXCreateElement"},5667:function(e,t,n){n.r(t),n.d(t,{frontMatter:function(){return u},contentTitle:function(){return l},metadata:function(){return c},toc:function(){return p},default:function(){return d}});var r=n(7462),o=n(3366),i=(n(7294),n(3905)),a=["components"],u={},l="Deployment",c={unversionedId:"guides/deployment",id:"guides/deployment",isDocsHomePage:!1,title:"Deployment",description:"Siren docker image can be found on Docker hub here. You can run the image with",source:"@site/docs/guides/deployment.md",sourceDirName:"guides",slug:"/guides/deployment",permalink:"/siren/docs/guides/deployment",editUrl:"https://github.com/odpf/siren/edit/master/docs/docs/guides/deployment.md",tags:[],version:"current",lastUpdatedBy:"Abduh",lastUpdatedAt:1658206269,formattedLastUpdatedAt:"7/19/2022",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Monitoring",permalink:"/siren/docs/guides/monitoring"},next:{title:"Troubleshooting",permalink:"/siren/docs/guides/troubleshooting"}},p=[{value:"Deploying to Kubernetes",id:"deploying-to-kubernetes",children:[]}],s={toc:p};function d(e){var t=e.components,n=(0,o.Z)(e,a);return(0,i.kt)("wrapper",(0,r.Z)({},s,n,{components:t,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"deployment"},"Deployment"),(0,i.kt)("p",null,"Siren docker image can be found on Docker hub ",(0,i.kt)("a",{parentName:"p",href:"https://hub.docker.com/r/odpf/siren"},"here"),". You can run the image with\nits dependencies."),(0,i.kt)("p",null,"The dependencies are:"),(0,i.kt)("ol",null,(0,i.kt)("li",{parentName:"ol"},"Postgres DB"),(0,i.kt)("li",{parentName:"ol"},"Cortex Ruler"),(0,i.kt)("li",{parentName:"ol"},"Cortex Alertmanager")),(0,i.kt)("p",null,"Make sure you have the instances running for them."),(0,i.kt)("h2",{id:"deploying-to-kubernetes"},"Deploying to Kubernetes"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},"Create a siren deployment using the helm chart available ",(0,i.kt)("a",{parentName:"li",href:"https://github.com/odpf/charts/tree/main/stable/siren"},"here"))))}d.isMDXComponent=!0}}]);