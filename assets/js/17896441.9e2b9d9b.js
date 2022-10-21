"use strict";(self.webpackChunksiren=self.webpackChunksiren||[]).push([[7918],{6742:function(e,t,n){n.d(t,{Z:function(){return v}});var a=n(3366),r=n(7294),i=n(3727),l=n(2263),o=n(3919),c=n(412),s=(0,r.createContext)({collectLink:function(){}}),d=n(4996),u=n(8780),m=["isNavLink","to","href","activeClassName","isActive","data-noBrokenLinkCheck","autoAddBaseUrl"];var v=function(e){var t,n,v=e.isNavLink,f=e.to,h=e.href,p=e.activeClassName,E=e.isActive,b=e["data-noBrokenLinkCheck"],g=e.autoAddBaseUrl,k=void 0===g||g,_=(0,a.Z)(e,m),Z=(0,l.Z)().siteConfig,N=Z.trailingSlash,U=Z.baseUrl,w=(0,d.C)().withBaseUrl,y=(0,r.useContext)(s),L=f||h,C=(0,o.Z)(L),T=null==L?void 0:L.replace("pathname://",""),A=void 0!==T?(n=T,k&&function(e){return e.startsWith("/")}(n)?w(n):n):void 0;A&&C&&(A=(0,u.applyTrailingSlash)(A,{trailingSlash:N,baseUrl:U}));var M,B=(0,r.useRef)(!1),S=v?i.OL:i.rU,O=c.Z.canUseIntersectionObserver;(0,r.useEffect)((function(){return!O&&C&&null!=A&&window.docusaurus.prefetch(A),function(){O&&M&&M.disconnect()}}),[A,O,C]);var x=null!==(t=null==A?void 0:A.startsWith("#"))&&void 0!==t&&t,I=!A||!C||x;return A&&C&&!x&&!b&&y.collectLink(A),I?r.createElement("a",Object.assign({href:A},L&&!C&&{target:"_blank",rel:"noopener noreferrer"},_)):r.createElement(S,Object.assign({},_,{onMouseEnter:function(){B.current||null==A||(window.docusaurus.preload(A),B.current=!0)},innerRef:function(e){var t,n;O&&e&&C&&(t=e,n=function(){null!=A&&window.docusaurus.prefetch(A)},(M=new window.IntersectionObserver((function(e){e.forEach((function(e){t===e.target&&(e.isIntersecting||e.intersectionRatio>0)&&(M.unobserve(t),M.disconnect(),n())}))}))).observe(t))},to:A||""},v&&{isActive:E,activeClassName:p}))}},3919:function(e,t,n){function a(e){return!0===/^(\w*:|\/\/)/.test(e)}function r(e){return void 0!==e&&!a(e)}n.d(t,{b:function(){return a},Z:function(){return r}})},4996:function(e,t,n){n.d(t,{C:function(){return i},Z:function(){return l}});var a=n(2263),r=n(3919);function i(){var e=(0,a.Z)().siteConfig,t=(e=void 0===e?{}:e).baseUrl,n=void 0===t?"/":t,i=e.url;return{withBaseUrl:function(e,t){return function(e,t,n,a){var i=void 0===a?{}:a,l=i.forcePrependBaseUrl,o=void 0!==l&&l,c=i.absolute,s=void 0!==c&&c;if(!n)return n;if(n.startsWith("#"))return n;if((0,r.b)(n))return n;if(o)return t+n;var d=n.startsWith(t)?n:t+n.replace(/^\//,"");return s?e+d:d}(i,n,e,t)}}}function l(e,t){return void 0===t&&(t={}),(0,i().withBaseUrl)(e,t)}},9362:function(e,t,n){n.r(t),n.d(t,{default:function(){return R}});var a=n(7294),r=n(6010),i=n(3783),l=n(6742),o=n(4973);var c=function(e){var t=e.metadata;return a.createElement("nav",{className:"pagination-nav docusaurus-mt-lg","aria-label":(0,o.I)({id:"theme.docs.paginator.navAriaLabel",message:"Docs pages navigation",description:"The ARIA label for the docs pagination"})},a.createElement("div",{className:"pagination-nav__item"},t.previous&&a.createElement(l.Z,{className:"pagination-nav__link",to:t.previous.permalink},a.createElement("div",{className:"pagination-nav__sublabel"},a.createElement(o.Z,{id:"theme.docs.paginator.previous",description:"The label used to navigate to the previous doc"},"Previous")),a.createElement("div",{className:"pagination-nav__label"},"\xab ",t.previous.title))),a.createElement("div",{className:"pagination-nav__item pagination-nav__item--next"},t.next&&a.createElement(l.Z,{className:"pagination-nav__link",to:t.next.permalink},a.createElement("div",{className:"pagination-nav__sublabel"},a.createElement(o.Z,{id:"theme.docs.paginator.next",description:"The label used to navigate to the next doc"},"Next")),a.createElement("div",{className:"pagination-nav__label"},t.next.title," \xbb"))))},s=n(2263),d=n(907),u=n(941);var m={unreleased:function(e){var t=e.siteTitle,n=e.versionMetadata;return a.createElement(o.Z,{id:"theme.docs.versions.unreleasedVersionLabel",description:"The label used to tell the user that he's browsing an unreleased doc version",values:{siteTitle:t,versionLabel:a.createElement("b",null,n.label)}},"This is unreleased documentation for {siteTitle} {versionLabel} version.")},unmaintained:function(e){var t=e.siteTitle,n=e.versionMetadata;return a.createElement(o.Z,{id:"theme.docs.versions.unmaintainedVersionLabel",description:"The label used to tell the user that he's browsing an unmaintained doc version",values:{siteTitle:t,versionLabel:a.createElement("b",null,n.label)}},"This is documentation for {siteTitle} {versionLabel}, which is no longer actively maintained.")}};function v(e){var t=m[e.versionMetadata.banner];return a.createElement(t,e)}function f(e){var t=e.versionLabel,n=e.to,r=e.onClick;return a.createElement(o.Z,{id:"theme.docs.versions.latestVersionSuggestionLabel",description:"The label used to tell the user to check the latest version",values:{versionLabel:t,latestVersionLink:a.createElement("b",null,a.createElement(l.Z,{to:n,onClick:r},a.createElement(o.Z,{id:"theme.docs.versions.latestVersionLinkLabel",description:"The label used for the latest version suggestion link label"},"latest version")))}},"For up-to-date documentation, see the {latestVersionLink} ({versionLabel}).")}function h(e){var t,n=e.versionMetadata,i=(0,s.Z)().siteConfig.title,l=(0,d.gA)({failfast:!0}).pluginId,o=(0,u.J)(l).savePreferredVersionName,c=(0,d.Jo)(l),m=c.latestDocSuggestion,h=c.latestVersionSuggestion,p=null!=m?m:(t=h).docs.find((function(e){return e.id===t.mainDocId}));return a.createElement("div",{className:(0,r.Z)(u.kM.docs.docVersionBanner,"alert alert--warning margin-bottom--md"),role:"alert"},a.createElement("div",null,a.createElement(v,{siteTitle:i,versionMetadata:n})),a.createElement("div",{className:"margin-top--md"},a.createElement(f,{versionLabel:h.label,to:p.path,onClick:function(){return o(h.name)}})))}var p=function(e){var t=e.versionMetadata;return t.banner?a.createElement(h,{versionMetadata:t}):a.createElement(a.Fragment,null)},E=n(1217);function b(e){var t=e.lastUpdatedAt,n=e.formattedLastUpdatedAt;return a.createElement(o.Z,{id:"theme.lastUpdated.atDate",description:"The words used to describe on which date a page has been last updated",values:{date:a.createElement("b",null,a.createElement("time",{dateTime:new Date(1e3*t).toISOString()},n))}}," on {date}")}function g(e){var t=e.lastUpdatedBy;return a.createElement(o.Z,{id:"theme.lastUpdated.byUser",description:"The words used to describe by who the page has been last updated",values:{user:a.createElement("b",null,t)}}," by {user}")}function k(e){var t=e.lastUpdatedAt,n=e.formattedLastUpdatedAt,r=e.lastUpdatedBy;return a.createElement("span",{className:u.kM.common.lastUpdated},a.createElement(o.Z,{id:"theme.lastUpdated.lastUpdatedAtBy",description:"The sentence used to display when a page has been last updated, and by who",values:{atDate:t&&n?a.createElement(b,{lastUpdatedAt:t,formattedLastUpdatedAt:n}):"",byUser:r?a.createElement(g,{lastUpdatedBy:r}):""}},"Last updated{atDate}{byUser}"),!1)}var _=n(6146),Z=n(7682),N="lastUpdated_mt2f";function U(e){return a.createElement("div",{className:(0,r.Z)(u.kM.docs.docFooterTagsRow,"row margin-bottom--sm")},a.createElement("div",{className:"col"},a.createElement(Z.Z,e)))}function w(e){var t=e.editUrl,n=e.lastUpdatedAt,i=e.lastUpdatedBy,l=e.formattedLastUpdatedAt;return a.createElement("div",{className:(0,r.Z)(u.kM.docs.docFooterEditMetaRow,"row")},a.createElement("div",{className:"col"},t&&a.createElement(_.Z,{editUrl:t})),a.createElement("div",{className:(0,r.Z)("col",N)},(n||i)&&a.createElement(k,{lastUpdatedAt:n,formattedLastUpdatedAt:l,lastUpdatedBy:i})))}function y(e){var t=e.content.metadata,n=t.editUrl,i=t.lastUpdatedAt,l=t.formattedLastUpdatedAt,o=t.lastUpdatedBy,c=t.tags,s=c.length>0,d=!!(n||i||o);return s||d?a.createElement("footer",{className:(0,r.Z)(u.kM.docs.docFooter,"docusaurus-mt-lg")},s&&a.createElement(U,{tags:c}),d&&a.createElement(w,{editUrl:n,lastUpdatedAt:i,lastUpdatedBy:o,formattedLastUpdatedAt:l})):a.createElement(a.Fragment,null)}var L=n(571),C="tocCollapsible_aw-L",T="tocCollapsibleButton_zr6a",A="tocCollapsibleContent_0dom",M="tocCollapsibleExpanded_FSiv";function B(e){var t,n=e.toc,i=e.className,l=(0,u.uR)({initialState:!0}),c=l.collapsed,s=l.toggleCollapsed;return a.createElement("div",{className:(0,r.Z)(C,(t={},t[M]=!c,t),i)},a.createElement("button",{type:"button",className:(0,r.Z)("clean-btn",T),onClick:s},a.createElement(o.Z,{id:"theme.TOCCollapsible.toggleButtonLabel",description:"The label used by the button on the collapsible TOC component"},"On this page")),a.createElement(u.zF,{lazy:!0,className:A,collapsed:c},a.createElement(L.r,{toc:n})))}var S=n(6159),O="docItemContainer_oiyr",x="docItemCol_zHA2",I="tocMobile_Tx6Y";function R(e){var t,n=e.content,l=e.versionMetadata,o=n.metadata,s=n.frontMatter,d=s.image,m=s.keywords,v=s.hide_title,f=s.hide_table_of_contents,h=o.description,b=o.title,g=!v&&void 0===n.contentTitle,k=(0,i.Z)(),_=!f&&n.toc&&n.toc.length>0,Z=_&&("desktop"===k||"ssr"===k);return a.createElement(a.Fragment,null,a.createElement(E.Z,{title:b,description:h,keywords:m,image:d}),a.createElement("div",{className:"row"},a.createElement("div",{className:(0,r.Z)("col",(t={},t[x]=!f,t))},a.createElement(p,{versionMetadata:l}),a.createElement("div",{className:O},a.createElement("article",null,l.badge&&a.createElement("span",{className:(0,r.Z)(u.kM.docs.docVersionBadge,"badge badge--secondary")},"Version: ",l.label),_&&a.createElement(B,{toc:n.toc,className:(0,r.Z)(u.kM.docs.docTocMobile,I)}),a.createElement("div",{className:(0,r.Z)(u.kM.docs.docMarkdown,"markdown")},g&&a.createElement(S.N,null,b),a.createElement(n,null)),a.createElement(y,e)),a.createElement(c,{metadata:o}))),Z&&a.createElement("div",{className:"col col--3"},a.createElement(L.Z,{toc:n.toc,className:u.kM.docs.docTocDesktop}))))}},6146:function(e,t,n){n.d(t,{Z:function(){return m}});var a=n(7294),r=n(4973),i=n(7462),l=n(3366),o=n(6010),c="iconEdit_mS5F",s=["className"],d=function(e){var t=e.className,n=(0,l.Z)(e,s);return a.createElement("svg",(0,i.Z)({fill:"currentColor",height:"20",width:"20",viewBox:"0 0 40 40",className:(0,o.Z)(c,t),"aria-hidden":"true"},n),a.createElement("g",null,a.createElement("path",{d:"m34.5 11.7l-3 3.1-6.3-6.3 3.1-3q0.5-0.5 1.2-0.5t1.1 0.5l3.9 3.9q0.5 0.4 0.5 1.1t-0.5 1.2z m-29.5 17.1l18.4-18.5 6.3 6.3-18.4 18.4h-6.3v-6.2z"})))},u=n(941);function m(e){var t=e.editUrl;return a.createElement("a",{href:t,target:"_blank",rel:"noreferrer noopener",className:u.kM.common.editThisPage},a.createElement(d,null),a.createElement(r.Z,{id:"theme.common.editThisPage",description:"The link label to edit the current page"},"Edit this page"))}},6159:function(e,t,n){n.d(t,{N:function(){return m},Z:function(){return v}});var a=n(3366),r=n(7462),i=n(7294),l=n(6010),o=n(4973),c=n(941),s="anchorWithStickyNavbar_y2LR",d="anchorWithHideOnScrollNavbar_3ly5",u=["id"],m=function(e){var t=Object.assign({},e);return i.createElement("header",null,i.createElement("h1",(0,r.Z)({},t,{id:void 0}),t.children))},v=function(e){return"h1"===e?m:(t=e,function(e){var n,r=e.id,m=(0,a.Z)(e,u),v=(0,c.LU)().navbar.hideOnScroll;return r?i.createElement(t,m,i.createElement("a",{"aria-hidden":"true",tabIndex:-1,className:(0,l.Z)("anchor","anchor__"+t,(n={},n[d]=v,n[s]=!v,n)),id:r}),m.children,i.createElement("a",{className:"hash-link",href:"#"+r,title:(0,o.I)({id:"theme.common.headingLinkTitle",message:"Direct link to heading",description:"Title for link to heading"})},"#")):i.createElement(t,m)});var t}},1217:function(e,t,n){n.d(t,{Z:function(){return o}});var a=n(7294),r=n(9105),i=n(941),l=n(4996);function o(e){var t=e.title,n=e.description,o=e.keywords,c=e.image,s=e.children,d=(0,i.pe)(t),u=(0,l.C)().withBaseUrl,m=c?u(c,{absolute:!0}):void 0;return a.createElement(r.Z,null,t&&a.createElement("title",null,d),t&&a.createElement("meta",{property:"og:title",content:d}),n&&a.createElement("meta",{name:"description",content:n}),n&&a.createElement("meta",{property:"og:description",content:n}),o&&a.createElement("meta",{name:"keywords",content:Array.isArray(o)?o.join(","):o}),m&&a.createElement("meta",{property:"og:image",content:m}),m&&a.createElement("meta",{name:"twitter:image",content:m}),s)}},571:function(e,t,n){n.d(t,{r:function(){return v},Z:function(){return f}});var a=n(7294),r=n(6010),i=n(941);function l(e){var t=e.getBoundingClientRect();return t.top===t.bottom?l(e.parentNode):t}function o(e){var t,n=e.anchorTopOffset,a=Array.from(document.querySelectorAll(".anchor.anchor__h2, .anchor.anchor__h3")),r=a.find((function(e){return l(e).top>=n}));return r?function(e){return e.top>0&&e.bottom<window.innerHeight/2}(l(r))?r:null!=(t=a[a.indexOf(r)-1])?t:null:a[a.length-1]}function c(){var e=(0,a.useRef)(0),t=(0,i.LU)().navbar.hideOnScroll;return(0,a.useEffect)((function(){e.current=t?0:document.querySelector(".navbar").clientHeight}),[t]),e}var s=function(e){var t=(0,a.useRef)(void 0),n=c();(0,a.useEffect)((function(){var a=e.linkClassName,r=e.linkActiveClassName;function i(){var e=function(e){return Array.from(document.getElementsByClassName(e))}(a),i=o({anchorTopOffset:n.current}),l=e.find((function(e){return i&&i.id===function(e){return decodeURIComponent(e.href.substring(e.href.indexOf("#")+1))}(e)}));e.forEach((function(e){!function(e,n){if(n){var a;t.current&&t.current!==e&&(null==(a=t.current)||a.classList.remove(r)),e.classList.add(r),t.current=e}else e.classList.remove(r)}(e,e===l)}))}return document.addEventListener("scroll",i),document.addEventListener("resize",i),i(),function(){document.removeEventListener("scroll",i),document.removeEventListener("resize",i)}}),[e,n])},d="tableOfContents_vrFS",u="table-of-contents__link",m={linkClassName:u,linkActiveClassName:"table-of-contents__link--active"};function v(e){var t=e.toc,n=e.isChild;return t.length?a.createElement("ul",{className:n?"":"table-of-contents table-of-contents__left-border"},t.map((function(e){return a.createElement("li",{key:e.id},a.createElement("a",{href:"#"+e.id,className:u,dangerouslySetInnerHTML:{__html:e.value}}),a.createElement(v,{isChild:!0,toc:e.children}))}))):null}var f=function(e){var t=e.toc;return s(m),a.createElement("div",{className:(0,r.Z)(d,"thin-scrollbar")},a.createElement(v,{toc:t}))}},7211:function(e,t,n){n.d(t,{Z:function(){return s}});var a=n(7294),r=n(6010),i=n(6742),l="tag_WK-t",o="tagRegular_LXbV",c="tagWithCount_S5Zl";var s=function(e){var t,n=e.permalink,s=e.name,d=e.count;return a.createElement(i.Z,{href:n,className:(0,r.Z)(l,(t={},t[o]=!d,t[c]=d,t))},s,d&&a.createElement("span",null,d))}},7682:function(e,t,n){n.d(t,{Z:function(){return s}});var a=n(7294),r=n(6010),i=n(4973),l=n(7211),o="tags_NBRY",c="tag_F03v";function s(e){var t=e.tags;return a.createElement(a.Fragment,null,a.createElement("b",null,a.createElement(i.Z,{id:"theme.tags.tagsListLabel",description:"The label alongside a tag list"},"Tags:")),a.createElement("ul",{className:(0,r.Z)(o,"padding--none","margin-left--sm")},t.map((function(e){var t=e.label,n=e.permalink;return a.createElement("li",{key:n,className:c},a.createElement(l.Z,{name:t,permalink:n}))}))))}},3783:function(e,t,n){var a=n(7294),r=n(412),i="desktop",l="mobile",o="ssr";function c(){return r.Z.canUseDOM?window.innerWidth>996?i:l:o}t.Z=function(){var e=(0,a.useState)((function(){return c()})),t=e[0],n=e[1];return(0,a.useEffect)((function(){function e(){n(c())}return window.addEventListener("resize",e),function(){window.removeEventListener("resize",e),clearTimeout(undefined)}}),[]),t}},8802:function(e,t){Object.defineProperty(t,"__esModule",{value:!0}),t.default=function(e,t){var n=t.trailingSlash,a=t.baseUrl;if(e.startsWith("#"))return e;if(void 0===n)return e;var r,i=e.split(/[#?]/)[0],l="/"===i||i===a?i:(r=i,n?function(e){return e.endsWith("/")?e:e+"/"}(r):function(e){return e.endsWith("/")?e.slice(0,-1):e}(r));return e.replace(i,l)}},8780:function(e,t,n){var a=this&&this.__importDefault||function(e){return e&&e.__esModule?e:{default:e}};Object.defineProperty(t,"__esModule",{value:!0}),t.uniq=t.applyTrailingSlash=void 0;var r=n(8802);Object.defineProperty(t,"applyTrailingSlash",{enumerable:!0,get:function(){return a(r).default}});var i=n(9964);Object.defineProperty(t,"uniq",{enumerable:!0,get:function(){return a(i).default}})},9964:function(e,t){Object.defineProperty(t,"__esModule",{value:!0}),t.default=function(e){return Array.from(new Set(e))}}}]);