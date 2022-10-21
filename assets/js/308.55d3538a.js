"use strict";(self.webpackChunksiren=self.webpackChunksiren||[]).push([[308],{6742:function(e,t,n){n.d(t,{Z:function(){return f}});var r=n(3366),a=n(7294),l=n(3727),o=n(2263),i=n(3919),c=n(412),s=(0,a.createContext)({collectLink:function(){}}),u=n(4996),m=n(8780),d=["isNavLink","to","href","activeClassName","isActive","data-noBrokenLinkCheck","autoAddBaseUrl"];var f=function(e){var t,n,f=e.isNavLink,v=e.to,h=e.href,b=e.activeClassName,g=e.isActive,E=e["data-noBrokenLinkCheck"],p=e.autoAddBaseUrl,Z=void 0===p||p,k=(0,r.Z)(e,d),w=(0,o.Z)().siteConfig,_=w.trailingSlash,N=w.baseUrl,y=(0,u.C)().withBaseUrl,C=(0,a.useContext)(s),L=v||h,I=(0,i.Z)(L),S=null==L?void 0:L.replace("pathname://",""),B=void 0!==S?(n=S,Z&&function(e){return e.startsWith("/")}(n)?y(n):n):void 0;B&&I&&(B=(0,m.applyTrailingSlash)(B,{trailingSlash:_,baseUrl:N}));var D,T=(0,a.useRef)(!1),A=f?l.OL:l.rU,U=c.Z.canUseIntersectionObserver;(0,a.useEffect)((function(){return!U&&I&&null!=B&&window.docusaurus.prefetch(B),function(){U&&D&&D.disconnect()}}),[B,U,I]);var M=null!==(t=null==B?void 0:B.startsWith("#"))&&void 0!==t&&t,x=!B||!I||M;return B&&I&&!M&&!E&&C.collectLink(B),x?a.createElement("a",Object.assign({href:B},L&&!I&&{target:"_blank",rel:"noopener noreferrer"},k)):a.createElement(A,Object.assign({},k,{onMouseEnter:function(){T.current||null==B||(window.docusaurus.preload(B),T.current=!0)},innerRef:function(e){var t,n;U&&e&&I&&(t=e,n=function(){null!=B&&window.docusaurus.prefetch(B)},(D=new window.IntersectionObserver((function(e){e.forEach((function(e){t===e.target&&(e.isIntersecting||e.intersectionRatio>0)&&(D.unobserve(t),D.disconnect(),n())}))}))).observe(t))},to:B||""},f&&{isActive:g,activeClassName:b}))}},1875:function(e,t){t.Z=function(){return null}},3919:function(e,t,n){function r(e){return!0===/^(\w*:|\/\/)/.test(e)}function a(e){return void 0!==e&&!r(e)}n.d(t,{b:function(){return r},Z:function(){return a}})},4996:function(e,t,n){n.d(t,{C:function(){return l},Z:function(){return o}});var r=n(2263),a=n(3919);function l(){var e=(0,r.Z)().siteConfig,t=(e=void 0===e?{}:e).baseUrl,n=void 0===t?"/":t,l=e.url;return{withBaseUrl:function(e,t){return function(e,t,n,r){var l=void 0===r?{}:r,o=l.forcePrependBaseUrl,i=void 0!==o&&o,c=l.absolute,s=void 0!==c&&c;if(!n)return n;if(n.startsWith("#"))return n;if((0,a.b)(n))return n;if(i)return t+n;var u=n.startsWith(t)?n:t+n.replace(/^\//,"");return s?e+u:u}(l,n,e,t)}}}function o(e,t){return void 0===t&&(t={}),(0,l().withBaseUrl)(e,t)}},8617:function(e,t,n){n.d(t,{Z:function(){return l}});var r=n(7294),a="iconExternalLink_wgqa",l=function(e){var t=e.width,n=void 0===t?13.5:t,l=e.height,o=void 0===l?13.5:l;return r.createElement("svg",{width:n,height:o,"aria-hidden":"true",viewBox:"0 0 24 24",className:a},r.createElement("path",{fill:"currentColor",d:"M21 13v10h-21v-19h12v2h-10v15h17v-8h2zm3-12h-10.988l4.035 4-6.977 7.07 2.828 2.828 6.977-7.07 4.125 4.172v-11z"}))}},308:function(e,t,n){n.d(t,{Z:function(){return Le}});var r=n(7294),a=n(6010),l=n(5977),o=n(4973),i=n(941),c="skipToContent_OuoZ";function s(e){e.setAttribute("tabindex","-1"),e.focus(),e.removeAttribute("tabindex")}var u=function(){var e=(0,r.useRef)(null),t=(0,l.k6)().action;return(0,i.SL)((function(n){var r=n.location;e.current&&!r.hash&&"PUSH"===t&&s(e.current)})),r.createElement("div",{ref:e},r.createElement("a",{href:"#",className:c,onClick:function(e){e.preventDefault();var t=document.querySelector("main:first-of-type")||document.querySelector(".main-wrapper");t&&s(t)}},r.createElement(o.Z,{id:"theme.common.skipToMainContent",description:"The skip to content label used for accessibility, allowing to rapidly navigate to main content with keyboard tab/enter navigation"},"Skip to main content")))},m=n(7462),d=n(3366),f=["width","height","className"];function v(e){var t=e.width,n=void 0===t?20:t,a=e.height,l=void 0===a?20:a,o=e.className,i=(0,d.Z)(e,f);return r.createElement("svg",(0,m.Z)({className:o,viewBox:"0 0 24 24",width:n,height:l,fill:"currentColor"},i),r.createElement("path",{d:"M24 20.188l-8.315-8.209 8.2-8.282-3.697-3.697-8.212 8.318-8.31-8.203-3.666 3.666 8.321 8.24-8.206 8.313 3.666 3.666 8.237-8.318 8.285 8.203z"}))}var h="announcementBar_axC9",b="announcementBarPlaceholder_xYHE",g="announcementBarClose_A3A1",E="announcementBarContent_6uhP";var p=function(){var e=(0,i.nT)(),t=e.isClosed,n=e.close,l=(0,i.LU)().announcementBar;if(!l)return null;var c=l.content,s=l.backgroundColor,u=l.textColor,m=l.isCloseable;return!c||m&&t?null:r.createElement("div",{className:h,style:{backgroundColor:s,color:u},role:"banner"},m&&r.createElement("div",{className:b}),r.createElement("div",{className:E,dangerouslySetInnerHTML:{__html:c}}),m?r.createElement("button",{type:"button",className:(0,a.Z)("clean-btn close",g),onClick:n,"aria-label":(0,o.I)({id:"theme.AnnouncementBar.closeButtonAriaLabel",message:"Close",description:"The ARIA label for close button of announcement bar"})},r.createElement(v,{width:14,height:14})):null)},Z=n(1875),k=n(2389),w={toggle:"toggle_iYfV"},_=function(e){var t=e.icon,n=e.style;return r.createElement("span",{className:(0,a.Z)(w.toggle,w.dark),style:n},t)},N=function(e){var t=e.icon,n=e.style;return r.createElement("span",{className:(0,a.Z)(w.toggle,w.light),style:n},t)},y=(0,r.memo)((function(e){var t=e.className,n=e.icons,l=e.checked,o=e.disabled,i=e.onChange,c=(0,r.useState)(l),s=c[0],u=c[1],m=(0,r.useState)(!1),d=m[0],f=m[1],v=(0,r.useRef)(null);return r.createElement("div",{className:(0,a.Z)("react-toggle",t,{"react-toggle--checked":s,"react-toggle--focus":d,"react-toggle--disabled":o})},r.createElement("div",{className:"react-toggle-track",role:"button",tabIndex:-1,onClick:function(){var e;return null==(e=v.current)?void 0:e.click()}},r.createElement("div",{className:"react-toggle-track-check"},n.checked),r.createElement("div",{className:"react-toggle-track-x"},n.unchecked),r.createElement("div",{className:"react-toggle-thumb"})),r.createElement("input",{ref:v,checked:s,type:"checkbox",className:"react-toggle-screenreader-only","aria-label":"Switch between dark and light mode",onChange:i,onClick:function(){return u(!s)},onFocus:function(){return f(!0)},onBlur:function(){return f(!1)},onKeyDown:function(e){var t;"Enter"===e.key&&(null==(t=v.current)||t.click())}}))}));function C(e){var t=(0,i.LU)().colorMode.switchConfig,n=t.darkIcon,a=t.darkIconStyle,l=t.lightIcon,o=t.lightIconStyle,c=(0,k.Z)();return r.createElement(y,(0,m.Z)({disabled:!c,icons:{checked:r.createElement(_,{icon:n,style:a}),unchecked:r.createElement(N,{icon:l,style:o})}},e))}var L=n(5350),I=n(7898),S=function(e){var t=(0,l.TH)(),n=(0,r.useState)(e),a=n[0],o=n[1],c=(0,r.useRef)(!1),s=(0,r.useState)(0),u=s[0],m=s[1],d=(0,r.useCallback)((function(e){null!==e&&m(e.getBoundingClientRect().height)}),[]);return(0,I.Z)((function(t,n){var r=t.scrollY,a=null==n?void 0:n.scrollY;if(e)if(r<u)o(!0);else{if(c.current)return c.current=!1,void o(!1);a&&0===r&&o(!0);var l=document.documentElement.scrollHeight-u,i=window.innerHeight;a&&r>=a?o(!1):r+i<l&&o(!0)}}),[u,c]),(0,i.SL)((function(t){e&&!t.location.hash&&o(!0)})),(0,r.useEffect)((function(){e&&t.hash&&(c.current=!0)}),[t.hash]),{navbarRef:d,isNavbarVisible:a}};var B=function(e){void 0===e&&(e=!0),(0,r.useEffect)((function(){return document.body.style.overflow=e?"hidden":"visible",function(){document.body.style.overflow="visible"}}),[e])},D=n(3783),T=n(907),A=n(7819),U=n(5537),M=["width","height","className"],x=function(e){var t=e.width,n=void 0===t?30:t,a=e.height,l=void 0===a?30:a,o=e.className,i=(0,d.Z)(e,M);return r.createElement("svg",(0,m.Z)({className:o,width:n,height:l,viewBox:"0 0 30 30","aria-hidden":"true"},i),r.createElement("path",{stroke:"currentColor",strokeLinecap:"round",strokeMiterlimit:"10",strokeWidth:"2",d:"M4 7h22M4 15h22M4 23h22"}))},R=["width","height","className"];function P(e){var t=e.width,n=void 0===t?20:t,a=e.height,l=void 0===a?20:a,o=e.className,i=(0,d.Z)(e,R);return r.createElement("svg",(0,m.Z)({className:o,viewBox:"0 0 413.348 413.348",width:n,height:l,fill:"currentColor"},i),r.createElement("path",{d:"m413.348 24.354-24.354-24.354-182.32 182.32-182.32-182.32-24.354 24.354 182.32 182.32-182.32 182.32 24.354 24.354 182.32-182.32 182.32 182.32 24.354-24.354-182.32-182.32z"}))}var O="toggle_2i4l",W="navbarHideable_RReh",H="navbarHidden_FBwS",V="navbarSidebarToggle_AVbO",z="navbarSidebarCloseSvg_+9jJ",j="right";function q(){return(0,i.LU)().navbar.items}function F(){var e=(0,i.LU)().colorMode.disableSwitch,t=(0,L.Z)(),n=t.isDarkTheme,a=t.setLightTheme,l=t.setDarkTheme;return{isDarkTheme:n,toggle:(0,r.useCallback)((function(e){return e.target.checked?l():a()}),[a,l]),disabled:e}}function G(e){var t=e.sidebarShown,n=e.toggleSidebar;B(t);var l=q(),c=F(),s=function(e){var t,n=e.sidebarShown,a=e.toggleSidebar,l=null==(t=(0,i.g8)())?void 0:t({toggleSidebar:a}),o=(0,i.D9)(l),c=(0,r.useState)((function(){return!1})),s=c[0],u=c[1];(0,r.useEffect)((function(){l&&!o&&u(!0)}),[l,o]);var m=!!l;return(0,r.useEffect)((function(){m?n||u(!0):u(!1)}),[n,m]),{shown:s,hide:(0,r.useCallback)((function(){u(!1)}),[]),content:l}}({sidebarShown:t,toggleSidebar:n});return r.createElement("div",{className:"navbar-sidebar"},r.createElement("div",{className:"navbar-sidebar__brand"},r.createElement(U.Z,{className:"navbar__brand",imageClassName:"navbar__logo",titleClassName:"navbar__title"}),!c.disabled&&r.createElement(C,{className:V,checked:c.isDarkTheme,onChange:c.toggle}),r.createElement("button",{type:"button",className:"clean-btn navbar-sidebar__close",onClick:n},r.createElement(P,{width:20,height:20,className:z}))),r.createElement("div",{className:(0,a.Z)("navbar-sidebar__items",{"navbar-sidebar__items--show-secondary":s.shown})},r.createElement("div",{className:"navbar-sidebar__item menu"},r.createElement("ul",{className:"menu__list"},l.map((function(e,t){return r.createElement(A.Z,(0,m.Z)({mobile:!0},e,{onClick:n,key:t}))})))),r.createElement("div",{className:"navbar-sidebar__item menu"},l.length>0&&r.createElement("button",{type:"button",className:"clean-btn navbar-sidebar__back",onClick:s.hide},r.createElement(o.Z,{id:"theme.navbar.mobileSidebarSecondaryMenu.backButtonLabel",description:"The label of the back button to return to main menu, inside the mobile navbar sidebar secondary menu (notably used to display the docs sidebar)"},"\u2190 Back to main menu")),s.content)))}var Y=function(){var e,t=(0,i.LU)().navbar,n=t.hideOnScroll,l=t.style,o=function(){var e=(0,D.Z)(),t="mobile"===e,n=(0,r.useState)(!1),a=n[0],l=n[1];(0,i.Rb)((function(){a&&l(!1)}));var o=(0,r.useCallback)((function(){l((function(e){return!e}))}),[]);return(0,r.useEffect)((function(){"desktop"===e&&l(!1)}),[e]),{shouldRender:t,toggle:o,shown:a}}(),c=F(),s=(0,T.gA)(),u=S(n),d=u.navbarRef,f=u.isNavbarVisible,v=q(),h=v.some((function(e){return"search"===e.type})),b=function(e){return{leftItems:e.filter((function(e){var t;return"left"===(null!=(t=e.position)?t:j)})),rightItems:e.filter((function(e){var t;return"right"===(null!=(t=e.position)?t:j)}))}}(v),g=b.leftItems,E=b.rightItems;return r.createElement("nav",{ref:d,className:(0,a.Z)("navbar","navbar--fixed-top",(e={"navbar--dark":"dark"===l,"navbar--primary":"primary"===l,"navbar-sidebar--show":o.shown},e[W]=n,e[H]=n&&!f,e))},r.createElement("div",{className:"navbar__inner"},r.createElement("div",{className:"navbar__items"},((null==v?void 0:v.length)>0||s)&&r.createElement("button",{"aria-label":"Navigation bar toggle",className:"navbar__toggle clean-btn",type:"button",tabIndex:0,onClick:o.toggle,onKeyDown:o.toggle},r.createElement(x,null)),r.createElement(U.Z,{className:"navbar__brand",imageClassName:"navbar__logo",titleClassName:"navbar__title"}),g.map((function(e,t){return r.createElement(A.Z,(0,m.Z)({},e,{key:t}))}))),r.createElement("div",{className:"navbar__items navbar__items--right"},E.map((function(e,t){return r.createElement(A.Z,(0,m.Z)({},e,{key:t}))})),!c.disabled&&r.createElement(C,{className:O,checked:c.isDarkTheme,onChange:c.toggle}),!h&&r.createElement(Z.Z,null))),r.createElement("div",{role:"presentation",className:"navbar-sidebar__backdrop",onClick:o.toggle}),o.shouldRender&&r.createElement(G,{sidebarShown:o.shown,toggleSidebar:o.toggle}))},J=n(6742),K=n(4996),Q=n(3919),X="footerLogoLink_SRtH",$=n(8465),ee=n(8617),te=["to","href","label","prependBaseUrlToHref"];function ne(e){var t=e.to,n=e.href,a=e.label,l=e.prependBaseUrlToHref,o=(0,d.Z)(e,te),i=(0,K.Z)(t),c=(0,K.Z)(n,{forcePrependBaseUrl:!0});return r.createElement(J.Z,(0,m.Z)({className:"footer__link-item"},n?{href:l?c:n}:{to:i},o),n&&!(0,Q.Z)(n)?r.createElement("span",null,a,r.createElement(ee.Z,null)):a)}var re=function(e){var t=e.sources,n=e.alt;return r.createElement($.Z,{className:"footer__logo",alt:n,sources:t})};var ae=function(){var e=(0,i.LU)().footer,t=e||{},n=t.copyright,l=t.links,o=void 0===l?[]:l,c=t.logo,s=void 0===c?{}:c,u={light:(0,K.Z)(s.src),dark:(0,K.Z)(s.srcDark||s.src)};return e?r.createElement("footer",{className:(0,a.Z)("footer",{"footer--dark":"dark"===e.style})},r.createElement("div",{className:"container"},o&&o.length>0&&r.createElement("div",{className:"row footer__links"},o.map((function(e,t){return r.createElement("div",{key:t,className:"col footer__col"},null!=e.title?r.createElement("div",{className:"footer__title"},e.title):null,null!=e.items&&Array.isArray(e.items)&&e.items.length>0?r.createElement("ul",{className:"footer__items"},e.items.map((function(e,t){return e.html?r.createElement("li",{key:t,className:"footer__item",dangerouslySetInnerHTML:{__html:e.html}}):r.createElement("li",{key:e.href||e.to,className:"footer__item"},r.createElement(ne,e))}))):null)}))),(s||n)&&r.createElement("div",{className:"footer__bottom text--center"},s&&(s.src||s.srcDark)&&r.createElement("div",{className:"margin-bottom--sm"},s.href?r.createElement(J.Z,{href:s.href,className:X},r.createElement(re,{alt:s.alt,sources:u})):r.createElement(re,{alt:s.alt,sources:u})),n?r.createElement("div",{className:"footer__copyright",dangerouslySetInnerHTML:{__html:n}}):null))):null},le=n(412),oe=(0,i.WA)("theme"),ie="light",ce="dark",se=function(e){return e===ce?ce:ie},ue=function(e){(0,i.WA)("theme").set(se(e))},me=function(){var e=(0,i.LU)().colorMode,t=e.defaultMode,n=e.disableSwitch,a=e.respectPrefersColorScheme,l=(0,r.useState)(function(e){return le.Z.canUseDOM?se(document.documentElement.getAttribute("data-theme")):se(e)}(t)),o=l[0],c=l[1],s=(0,r.useCallback)((function(){c(ie),ue(ie)}),[]),u=(0,r.useCallback)((function(){c(ce),ue(ce)}),[]);return(0,r.useEffect)((function(){document.documentElement.setAttribute("data-theme",se(o))}),[o]),(0,r.useEffect)((function(){if(!n)try{var e=oe.get();null!==e&&c(se(e))}catch(t){console.error(t)}}),[c]),(0,r.useEffect)((function(){n&&!a||window.matchMedia("(prefers-color-scheme: dark)").addListener((function(e){var t=e.matches;c(t?ce:ie)}))}),[]),{isDarkTheme:o===ce,setLightTheme:s,setDarkTheme:u}},de=n(2924);var fe=function(e){var t=me(),n=t.isDarkTheme,a=t.setLightTheme,l=t.setDarkTheme;return r.createElement(de.Z.Provider,{value:{isDarkTheme:n,setLightTheme:a,setDarkTheme:l}},e.children)},ve="docusaurus.tab.",he=function(){var e=(0,r.useState)({}),t=e[0],n=e[1],a=(0,r.useCallback)((function(e,t){(0,i.WA)("docusaurus.tab."+e).set(t)}),[]);return(0,r.useEffect)((function(){try{var e={};(0,i._f)().forEach((function(t){if(t.startsWith(ve)){var n=t.substring(ve.length);e[n]=(0,i.WA)(t).get()}})),n(e)}catch(t){console.error(t)}}),[]),{tabGroupChoices:t,setTabGroupChoices:function(e,t){n((function(n){var r;return Object.assign({},n,((r={})[e]=t,r))})),a(e,t)}}},be=n(9443);var ge=function(e){var t=he(),n=t.tabGroupChoices,a=t.setTabGroupChoices;return r.createElement(be.Z.Provider,{value:{tabGroupChoices:n,setTabGroupChoices:a}},e.children)};function Ee(e){var t=e.children;return r.createElement(fe,null,r.createElement(i.pl,null,r.createElement(ge,null,r.createElement(i.L5,null,r.createElement(i.Cn,null,t)))))}var pe=n(9105),Ze=n(2263);function ke(e){var t=e.locale,n=e.version,a=e.tag;return r.createElement(pe.Z,null,t&&r.createElement("meta",{name:"docusaurus_locale",content:t}),n&&r.createElement("meta",{name:"docusaurus_version",content:n}),a&&r.createElement("meta",{name:"docusaurus_tag",content:a}))}var we=n(1217);function _e(){var e=(0,Ze.Z)().i18n,t=e.defaultLocale,n=e.locales,a=(0,i.l5)();return r.createElement(pe.Z,null,n.map((function(e){return r.createElement("link",{key:e,rel:"alternate",href:a.createUrl({locale:e,fullyQualified:!0}),hrefLang:e})})),r.createElement("link",{rel:"alternate",href:a.createUrl({locale:t,fullyQualified:!0}),hrefLang:"x-default"}))}function Ne(e){var t=e.permalink,n=(0,Ze.Z)().siteConfig.url,a=function(){var e=(0,Ze.Z)().siteConfig.url,t=(0,l.TH)().pathname;return e+(0,K.Z)(t)}(),o=t?""+n+t:a;return r.createElement(pe.Z,null,r.createElement("meta",{property:"og:url",content:o}),r.createElement("link",{rel:"canonical",href:o}))}function ye(e){var t=(0,Ze.Z)(),n=t.siteConfig.favicon,a=t.i18n,l=a.currentLocale,o=a.localeConfigs,c=(0,i.LU)(),s=c.metadatas,u=c.image,d=e.title,f=e.description,v=e.image,h=e.keywords,b=e.searchMetadatas,g=(0,K.Z)(n),E=(0,i.pe)(d),p=l,Z=o[l].direction;return r.createElement(r.Fragment,null,r.createElement(pe.Z,null,r.createElement("html",{lang:p,dir:Z}),n&&r.createElement("link",{rel:"shortcut icon",href:g}),r.createElement("title",null,E),r.createElement("meta",{property:"og:title",content:E}),r.createElement("meta",{name:"twitter:card",content:"summary_large_image"})),u&&r.createElement(we.Z,{image:u}),v&&r.createElement(we.Z,{image:v}),r.createElement(we.Z,{description:f,keywords:h}),r.createElement(Ne,null),r.createElement(_e,null),r.createElement(ke,(0,m.Z)({tag:i.HX,locale:l},b)),r.createElement(pe.Z,null,s.map((function(e,t){return r.createElement("meta",(0,m.Z)({key:"metadata_"+t},e))}))))}var Ce=function(){(0,r.useEffect)((function(){var e="navigation-with-keyboard";function t(t){"keydown"===t.type&&"Tab"===t.key&&document.body.classList.add(e),"mousedown"===t.type&&document.body.classList.remove(e)}return document.addEventListener("keydown",t),document.addEventListener("mousedown",t),function(){document.body.classList.remove(e),document.removeEventListener("keydown",t),document.removeEventListener("mousedown",t)}}),[])};var Le=function(e){var t=e.children,n=e.noFooter,l=e.wrapperClassName,o=e.pageClassName;return Ce(),r.createElement(Ee,null,r.createElement(ye,e),r.createElement(u,null),r.createElement(p,null),r.createElement(Y,null),r.createElement("div",{className:(0,a.Z)(i.kM.wrapper.main,l,o)},t),!n&&r.createElement(ae,null))}},5537:function(e,t,n){var r=n(7462),a=n(3366),l=n(7294),o=n(6742),i=n(8465),c=n(4996),s=n(2263),u=n(941),m=["imageClassName","titleClassName"];t.Z=function(e){var t=(0,s.Z)().siteConfig.title,n=(0,u.LU)().navbar,d=n.title,f=n.logo,v=void 0===f?{src:""}:f,h=e.imageClassName,b=e.titleClassName,g=(0,a.Z)(e,m),E=(0,c.Z)(v.href||"/"),p={light:(0,c.Z)(v.src),dark:(0,c.Z)(v.srcDark||v.src)};return l.createElement(o.Z,(0,r.Z)({to:E},g,v.target&&{target:v.target}),v.src&&l.createElement(i.Z,{className:h,sources:p,alt:v.alt||d||t}),null!=d&&l.createElement("b",{className:b},d))}},5525:function(e,t,n){n.d(t,{O:function(){return b}});var r=n(7462),a=n(3366),l=n(7294),o=n(6010),i=n(6742),c=n(4996),s=n(8617),u=n(3919),m=n(7819),d=["activeBasePath","activeBaseRegex","to","href","label","activeClassName","prependBaseUrlToHref"],f=["className","isDropdownItem"],v=["className","isDropdownItem"],h=["mobile","position"];function b(e){var t,n=e.activeBasePath,o=e.activeBaseRegex,m=e.to,f=e.href,v=e.label,h=e.activeClassName,b=void 0===h?"":h,g=e.prependBaseUrlToHref,E=(0,a.Z)(e,d),p=(0,c.Z)(m),Z=(0,c.Z)(n),k=(0,c.Z)(f,{forcePrependBaseUrl:!0}),w=v&&f&&!(0,u.Z)(f),_="dropdown__link--active"===b;return l.createElement(i.Z,(0,r.Z)({},f?{href:g?k:f}:Object.assign({isNavLink:!0,activeClassName:null!=(t=E.className)&&t.includes(b)?"":b,to:p},n||o?{isActive:function(e,t){return o?new RegExp(o).test(t.pathname):t.pathname.startsWith(Z)}}:null),E),w?l.createElement("span",null,v,l.createElement(s.Z,_&&{width:12,height:12})):v)}function g(e){var t=e.className,n=e.isDropdownItem,i=void 0!==n&&n,c=(0,a.Z)(e,f),s=l.createElement(b,(0,r.Z)({className:(0,o.Z)(i?"dropdown__link":"navbar__item navbar__link",t)},c));return i?l.createElement("li",null,s):s}function E(e){var t=e.className,n=(e.isDropdownItem,(0,a.Z)(e,v));return l.createElement("li",{className:"menu__list-item"},l.createElement(b,(0,r.Z)({className:(0,o.Z)("menu__link",t)},n)))}t.Z=function(e){var t,n=e.mobile,o=void 0!==n&&n,i=(e.position,(0,a.Z)(e,h)),c=o?E:g;return l.createElement(c,(0,r.Z)({},i,{activeClassName:null!=(t=i.activeClassName)?t:(0,m.E)(o)}))}},6400:function(e,t,n){n.d(t,{Z:function(){return f}});var r=n(7462),a=n(3366),l=n(7294),o=n(5525),i=n(907),c=n(6010),s=n(7819),u=n(941),m=n(8780),d=["docId","label","docsPluginId"];function f(e){var t,n=e.docId,f=e.label,v=e.docsPluginId,h=(0,a.Z)(e,d),b=(0,i.Iw)(v),g=b.activeVersion,E=b.activeDoc,p=(0,u.J)(v).preferredVersion,Z=(0,i.yW)(v),k=function(e,t){var n=e.flatMap((function(e){return e.docs})),r=n.find((function(e){return e.id===t}));if(!r){var a=n.map((function(e){return e.id})).join("\n- ");throw new Error("DocNavbarItem: couldn't find any doc with id \""+t+'" in version'+(e.length?"s":"")+" "+e.map((function(e){return e.name})).join(", ")+'".\nAvailable doc ids are:\n- '+a)}return r}((0,m.uniq)([g,p,Z].filter(Boolean)),n),w=(0,s.E)(h.mobile);return l.createElement(o.Z,(0,r.Z)({exact:!0},h,{className:(0,c.Z)(h.className,(t={},t[w]=(null==E?void 0:E.sidebar)&&E.sidebar===k.sidebar,t)),activeClassName:w,label:null!=f?f:k.id,to:k.path}))}},9308:function(e,t,n){n.d(t,{Z:function(){return f}});var r=n(7462),a=n(3366),l=n(7294),o=n(5525),i=n(3154),c=n(907),s=n(941),u=n(4973),m=["mobile","docsPluginId","dropdownActiveClassDisabled","dropdownItemsBefore","dropdownItemsAfter"],d=function(e){return e.docs.find((function(t){return t.id===e.mainDocId}))};function f(e){var t,n,f=e.mobile,v=e.docsPluginId,h=e.dropdownActiveClassDisabled,b=e.dropdownItemsBefore,g=e.dropdownItemsAfter,E=(0,a.Z)(e,m),p=(0,c.Iw)(v),Z=(0,c.gB)(v),k=(0,c.yW)(v),w=(0,s.J)(v),_=w.preferredVersion,N=w.savePreferredVersionName;var y,C=(y=Z.map((function(e){var t=(null==p?void 0:p.alternateDocVersions[e.name])||d(e);return{isNavLink:!0,label:e.label,to:t.path,isActive:function(){return e===(null==p?void 0:p.activeVersion)},onClick:function(){N(e.name)}}})),[].concat(b,y,g)),L=null!=(t=null!=(n=p.activeVersion)?n:_)?t:k,I=f&&C?(0,u.I)({id:"theme.navbar.mobileVersionsDropdown.label",message:"Versions",description:"The label for the navbar versions dropdown on mobile view"}):L.label,S=f&&C?void 0:d(L).path;return C.length<=1?l.createElement(o.Z,(0,r.Z)({},E,{mobile:f,label:I,to:S,isActive:h?function(){return!1}:void 0})):l.createElement(i.Z,(0,r.Z)({},E,{mobile:f,label:I,to:S,items:C,isActive:h?function(){return!1}:void 0}))}},7250:function(e,t,n){n.d(t,{Z:function(){return u}});var r=n(7462),a=n(3366),l=n(7294),o=n(5525),i=n(907),c=n(941),s=["label","to","docsPluginId"];function u(e){var t,n=e.label,u=e.to,m=e.docsPluginId,d=(0,a.Z)(e,s),f=(0,i.zu)(m),v=(0,c.J)(m).preferredVersion,h=(0,i.yW)(m),b=null!=(t=null!=f?f:v)?t:h,g=null!=n?n:b.label,E=null!=u?u:function(e){return e.docs.find((function(t){return t.id===e.mainDocId}))}(b).path;return l.createElement(o.Z,(0,r.Z)({},d,{label:g,to:E}))}},3154:function(e,t,n){var r=n(7462),a=n(3366),l=n(7294),o=n(6010),i=n(941),c=n(5525),s=n(7819),u=["items","position","className"],m=["items","className","position"],d=["mobile"];function f(e,t){return e.some((function(e){return function(e,t){return!!(0,i.Mg)(e.to,t)||!(!e.activeBaseRegex||!new RegExp(e.activeBaseRegex).test(t))||!(!e.activeBasePath||!t.startsWith(e.activeBasePath))}(e,t)}))}function v(e){var t,n=e.items,i=e.position,m=e.className,d=(0,a.Z)(e,u),f=(0,l.useRef)(null),v=(0,l.useRef)(null),h=(0,l.useState)(!1),b=h[0],g=h[1];return(0,l.useEffect)((function(){var e=function(e){f.current&&!f.current.contains(e.target)&&g(!1)};return document.addEventListener("mousedown",e),document.addEventListener("touchstart",e),function(){document.removeEventListener("mousedown",e),document.removeEventListener("touchstart",e)}}),[f]),l.createElement("div",{ref:f,className:(0,o.Z)("navbar__item","dropdown","dropdown--hoverable",{"dropdown--right":"right"===i,"dropdown--show":b})},l.createElement(c.O,(0,r.Z)({className:(0,o.Z)("navbar__link",m)},d,{onClick:d.to?void 0:function(e){return e.preventDefault()},onKeyDown:function(e){"Enter"===e.key&&(e.preventDefault(),g(!b))}}),null!=(t=d.children)?t:d.label),l.createElement("ul",{ref:v,className:"dropdown__menu"},n.map((function(e,t){return l.createElement(s.Z,(0,r.Z)({isDropdownItem:!0,onKeyDown:function(e){if(t===n.length-1&&"Tab"===e.key){e.preventDefault(),g(!1);var r=f.current.nextElementSibling;r&&r.focus()}},activeClassName:"dropdown__link--active"},e,{key:t}))}))))}function h(e){var t,n=e.items,u=e.className,d=(e.position,(0,a.Z)(e,m)),v=(0,i.be)(),h=f(n,v),b=(0,i.uR)({initialState:function(){return!h}}),g=b.collapsed,E=b.toggleCollapsed,p=b.setCollapsed;return(0,l.useEffect)((function(){h&&p(!h)}),[v,h]),l.createElement("li",{className:(0,o.Z)("menu__list-item",{"menu__list-item--collapsed":g})},l.createElement(c.O,(0,r.Z)({role:"button",className:(0,o.Z)("menu__link menu__link--sublist",u)},d,{onClick:function(e){e.preventDefault(),E()}}),null!=(t=d.children)?t:d.label),l.createElement(i.zF,{lazy:!0,as:"ul",className:"menu__list",collapsed:g},n.map((function(e,t){return l.createElement(s.Z,(0,r.Z)({mobile:!0,isDropdownItem:!0,onClick:d.onClick,activeClassName:"menu__link--active"},e,{key:t}))}))))}t.Z=function(e){var t=e.mobile,n=void 0!==t&&t,r=(0,a.Z)(e,d),o=n?h:v;return l.createElement(o,r)}},7819:function(e,t,n){n.d(t,{Z:function(){return Z},E:function(){return p}});var r=n(3366),a=n(7294),l=n(5525),o=n(3154),i=n(7462),c=["width","height"],s=function(e){var t=e.width,n=void 0===t?20:t,l=e.height,o=void 0===l?20:l,s=(0,r.Z)(e,c);return a.createElement("svg",(0,i.Z)({viewBox:"0 0 20 20",width:n,height:o,"aria-hidden":"true"},s),a.createElement("path",{fill:"currentColor",d:"M19.753 10.909c-.624-1.707-2.366-2.726-4.661-2.726-.09 0-.176.002-.262.006l-.016-2.063 3.525-.607c.115-.019.133-.119.109-.231-.023-.111-.167-.883-.188-.976-.027-.131-.102-.127-.207-.109-.104.018-3.25.461-3.25.461l-.013-2.078c-.001-.125-.069-.158-.194-.156l-1.025.016c-.105.002-.164.049-.162.148l.033 2.307s-3.061.527-3.144.543c-.084.014-.17.053-.151.143.019.09.19 1.094.208 1.172.018.08.072.129.188.107l2.924-.504.035 2.018c-1.077.281-1.801.824-2.256 1.303-.768.807-1.207 1.887-1.207 2.963 0 1.586.971 2.529 2.328 2.695 3.162.387 5.119-3.06 5.769-4.715 1.097 1.506.256 4.354-2.094 5.98-.043.029-.098.129-.033.207l.619.756c.08.096.206.059.256.023 2.51-1.73 3.661-4.515 2.869-6.683zm-7.386 3.188c-.966-.121-.944-.914-.944-1.453 0-.773.327-1.58.876-2.156a3.21 3.21 0 011.229-.799l.082 4.277a2.773 2.773 0 01-1.243.131zm2.427-.553l.046-4.109c.084-.004.166-.01.252-.01.773 0 1.494.145 1.885.361.391.217-1.023 2.713-2.183 3.758zm-8.95-7.668a.196.196 0 00-.196-.145h-1.95a.194.194 0 00-.194.144L.008 16.916c-.017.051-.011.076.062.076h1.733c.075 0 .099-.023.114-.072l1.008-3.318h3.496l1.008 3.318c.016.049.039.072.113.072h1.734c.072 0 .078-.025.062-.076-.014-.05-3.083-9.741-3.494-11.04zm-2.618 6.318l1.447-5.25 1.447 5.25H3.226z"}))},u=n(2263),m=n(941),d="iconLanguage_EbrZ",f=["mobile","dropdownItemsBefore","dropdownItemsAfter"];function v(e){var t=e.mobile,n=e.dropdownItemsBefore,l=e.dropdownItemsAfter,c=(0,r.Z)(e,f),v=(0,u.Z)().i18n,h=v.currentLocale,b=v.locales,g=v.localeConfigs,E=(0,m.l5)();function p(e){return g[e].label}var Z=b.map((function(e){var t="pathname://"+E.createUrl({locale:e,fullyQualified:!1});return{isNavLink:!0,label:p(e),to:t,target:"_self",autoAddBaseUrl:!1,className:e===h?"dropdown__link--active":"",style:{textTransform:"capitalize"}}})),k=[].concat(n,Z,l),w=t?"Languages":p(h);return a.createElement(o.Z,(0,i.Z)({},c,{href:"#",mobile:t,label:a.createElement("span",null,a.createElement(s,{className:d}),a.createElement("span",null,w)),items:k}))}var h=n(1875);function b(e){return e.mobile?null:a.createElement(h.Z,null)}var g=["type"],E={default:function(){return l.Z},localeDropdown:function(){return v},search:function(){return b},dropdown:function(){return o.Z},docsVersion:function(){return n(7250).Z},docsVersionDropdown:function(){return n(9308).Z},doc:function(){return n(6400).Z}};var p=function(e){return e?"menu__link--active":"navbar__link--active"};function Z(e){var t=e.type,n=(0,r.Z)(e,g),l=function(e,t){return e&&"default"!==e?e:t?"dropdown":"default"}(t,void 0!==n.items),o=function(e){var t=E[e];if(!t)throw new Error('No NavbarItem component found for type "'+e+'".');return t()}(l);return a.createElement(o,n)}},1217:function(e,t,n){n.d(t,{Z:function(){return i}});var r=n(7294),a=n(9105),l=n(941),o=n(4996);function i(e){var t=e.title,n=e.description,i=e.keywords,c=e.image,s=e.children,u=(0,l.pe)(t),m=(0,o.C)().withBaseUrl,d=c?m(c,{absolute:!0}):void 0;return r.createElement(a.Z,null,t&&r.createElement("title",null,u),t&&r.createElement("meta",{property:"og:title",content:u}),n&&r.createElement("meta",{name:"description",content:n}),n&&r.createElement("meta",{property:"og:description",content:n}),i&&r.createElement("meta",{name:"keywords",content:Array.isArray(i)?i.join(","):i}),d&&r.createElement("meta",{property:"og:image",content:d}),d&&r.createElement("meta",{name:"twitter:image",content:d}),s)}},8465:function(e,t,n){n.d(t,{Z:function(){return m}});var r=n(7462),a=n(3366),l=n(7294),o=n(6010),i=n(2389),c=n(5350),s={themedImage:"themedImage_TMUO","themedImage--light":"themedImage--light_4Vu1","themedImage--dark":"themedImage--dark_uzRr"},u=["sources","className","alt"],m=function(e){var t=(0,i.Z)(),n=(0,c.Z)().isDarkTheme,m=e.sources,d=e.className,f=e.alt,v=void 0===f?"":f,h=(0,a.Z)(e,u),b=t?n?["dark"]:["light"]:["light","dark"];return l.createElement(l.Fragment,null,b.map((function(e){return l.createElement("img",(0,r.Z)({key:e,src:m[e],alt:v,className:(0,o.Z)(s.themedImage,s["themedImage--"+e],d)},h))})))}},7898:function(e,t,n){var r=n(7294),a=n(412),l=function(){return a.Z.canUseDOM?{scrollX:window.pageXOffset,scrollY:window.pageYOffset}:null};t.Z=function(e,t){void 0===t&&(t=[]);var n=(0,r.useRef)(l()),a=function(){var t=l();e&&e(t,n.current),n.current=t};(0,r.useEffect)((function(){var e={passive:!0};return a(),window.addEventListener("scroll",a,e),function(){return window.removeEventListener("scroll",a,e)}}),t)}},3783:function(e,t,n){var r=n(7294),a=n(412),l="desktop",o="mobile",i="ssr";function c(){return a.Z.canUseDOM?window.innerWidth>996?l:o:i}t.Z=function(){var e=(0,r.useState)((function(){return c()})),t=e[0],n=e[1];return(0,r.useEffect)((function(){function e(){n(c())}return window.addEventListener("resize",e),function(){window.removeEventListener("resize",e),clearTimeout(undefined)}}),[]),t}},8802:function(e,t){Object.defineProperty(t,"__esModule",{value:!0}),t.default=function(e,t){var n=t.trailingSlash,r=t.baseUrl;if(e.startsWith("#"))return e;if(void 0===n)return e;var a,l=e.split(/[#?]/)[0],o="/"===l||l===r?l:(a=l,n?function(e){return e.endsWith("/")?e:e+"/"}(a):function(e){return e.endsWith("/")?e.slice(0,-1):e}(a));return e.replace(l,o)}},8780:function(e,t,n){var r=this&&this.__importDefault||function(e){return e&&e.__esModule?e:{default:e}};Object.defineProperty(t,"__esModule",{value:!0}),t.uniq=t.applyTrailingSlash=void 0;var a=n(8802);Object.defineProperty(t,"applyTrailingSlash",{enumerable:!0,get:function(){return r(a).default}});var l=n(9964);Object.defineProperty(t,"uniq",{enumerable:!0,get:function(){return r(l).default}})},9964:function(e,t){Object.defineProperty(t,"__esModule",{value:!0}),t.default=function(e){return Array.from(new Set(e))}}}]);