"use strict";(self.webpackChunksiren=self.webpackChunksiren||[]).push([[242],{3905:function(e,n,t){t.d(n,{Zo:function(){return u},kt:function(){return d}});var r=t(7294);function i(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function o(e,n){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);n&&(r=r.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),t.push.apply(t,r)}return t}function a(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{};n%2?o(Object(t),!0).forEach((function(n){i(e,n,t[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):o(Object(t)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(t,n))}))}return e}function c(e,n){if(null==e)return{};var t,r,i=function(e,n){if(null==e)return{};var t,r,i={},o=Object.keys(e);for(r=0;r<o.length;r++)t=o[r],n.indexOf(t)>=0||(i[t]=e[t]);return i}(e,n);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)t=o[r],n.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(i[t]=e[t])}return i}var l=r.createContext({}),s=function(e){var n=r.useContext(l),t=n;return e&&(t="function"==typeof e?e(n):a(a({},n),e)),t},u=function(e){var n=s(e.components);return r.createElement(l.Provider,{value:n},e.children)},f={inlineCode:"code",wrapper:function(e){var n=e.children;return r.createElement(r.Fragment,{},n)}},p=r.forwardRef((function(e,n){var t=e.components,i=e.mdxType,o=e.originalType,l=e.parentName,u=c(e,["components","mdxType","originalType","parentName"]),p=s(t),d=i,m=p["".concat(l,".").concat(d)]||p[d]||f[d]||o;return t?r.createElement(m,a(a({ref:n},u),{},{components:t})):r.createElement(m,a({ref:n},u))}));function d(e,n){var t=arguments,i=n&&n.mdxType;if("string"==typeof e||i){var o=t.length,a=new Array(o);a[0]=p;var c={};for(var l in n)hasOwnProperty.call(n,l)&&(c[l]=n[l]);c.originalType=e,c.mdxType="string"==typeof e?e:i,a[1]=c;for(var s=2;s<o;s++)a[s]=t[s];return r.createElement.apply(null,a)}return r.createElement.apply(null,t)}p.displayName="MDXCreateElement"},1925:function(e,n,t){t.r(n),t.d(n,{assets:function(){return l},contentTitle:function(){return a},default:function(){return f},frontMatter:function(){return o},metadata:function(){return c},toc:function(){return s}});var r=t(3117),i=(t(7294),t(3905));const o={},a="Client Configuration",c={unversionedId:"reference/client_configuration",id:"reference/client_configuration",title:"Client Configuration",description:"When using siren client CLI, sometimes there are client-specifi flags that are required to be passed e.g. --host so you are calling Siren like this.",source:"@site/docs/reference/client_configuration.md",sourceDirName:"reference",slug:"/reference/client_configuration",permalink:"/siren/docs/reference/client_configuration",draft:!1,editUrl:"https://github.com/odpf/siren/edit/master/docs/docs/reference/client_configuration.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Server Configuration",permalink:"/siren/docs/reference/server_configuration"},next:{title:"Receiver",permalink:"/siren/docs/reference/receiver"}},l={},s=[],u={toc:s};function f(e){let{components:n,...t}=e;return(0,i.kt)("wrapper",(0,r.Z)({},u,t,{components:n,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"client-configuration"},"Client Configuration"),(0,i.kt)("p",null,"When using ",(0,i.kt)("inlineCode",{parentName:"p"},"siren")," client CLI, sometimes there are client-specifi flags that are required to be passed e.g. ",(0,i.kt)("inlineCode",{parentName:"p"},"--host")," so you are calling Siren like this."),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-bash"},"siren receiver list --host localhost:8080\n")),(0,i.kt)("p",null,"Siren client CLI could use a client config so you don't need to pass client-required flags e.g. ",(0,i.kt)("inlineCode",{parentName:"p"},"--host")," everytime you run ",(0,i.kt)("inlineCode",{parentName:"p"},"siren")," command. For Siren client CLI, here is the required config to interact to Siren server."),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-yaml"},"host: localhost:8080\n")),(0,i.kt)("p",null,"You could easily generate client config by running this command:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-bash"},"siren config init\n")),(0,i.kt)("p",null,"This will create (if not exists) a config file ",(0,i.kt)("inlineCode",{parentName:"p"},"${HOME}/.config/odpf/siren.yaml")," with default values. You can modify the value as you wish."))}f.isMDXComponent=!0}}]);