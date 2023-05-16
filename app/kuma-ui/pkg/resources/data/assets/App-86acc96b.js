import{d as f,u as K,a as Q,c as y,o,b as l,e as s,f as h,w as _,g as v,h as n,i as t,j as u,t as A,r as $,k as B,l as V,n as q,p as O,m as E,P as j,F as z,q as F,s as D,C as R,v as x,x as H,y as Z,T as G}from"./index-0fbacd76.js";import{d as J,c as T,q as X,T as N,Z as W,z as ee,o as te}from"./kongponents.es-c741eab8.js";import{u as k}from"./store-5ee7e2bf.js";import{u as se,a as S}from"./index-bc240554.js";import{_ as U}from"./_plugin-vue_export-helper-c27b6911.js";import{A as oe,a as ne}from"./AccordionList-1e0cd60d.js";import{u as ae,a as ie}from"./index-93c6e791.js";import"./datadogLogEvents-302eea7b.js";import"./DoughnutChart-46244ceb.js";const ce=f({__name:"AppBreadcrumbs",setup(r){const e=K(),i=Q(),c=k(),d=y(()=>e.matched.filter(a=>a.meta.isBreadcrumb===!0).map(a=>{try{const m=i.resolve(a);if(typeof m.path=="string"&&m.path!==""){const b={path:m.path};return{matchedRoute:a,to:b}}else return null}catch{return null}}).filter(p).map(({matchedRoute:a,to:m})=>{const b=M(a,e);return{to:m,title:b,text:b}}));function p(a){return a!==null}function M(a,m){return typeof a.meta.getBreadcrumbTitle=="function"?a.meta.getBreadcrumbTitle(m,c):a.meta.breadcrumbTitleParam&&m.params[a.meta.breadcrumbTitleParam]?m.params[a.meta.breadcrumbTitleParam]:a.meta.title}return(a,m)=>d.value.length>0?(o(),l(s(J),{key:0,items:d.value},null,8,["items"])):h("",!0)}}),re=t("p",null,"Unable to reach the API",-1),le={key:0},_e=f({__name:"AppErrorMessage",setup(r){const e=se();return(i,c)=>(o(),l(s(X),{class:"global-api-status empty-state--wide-content empty-state--compact","cta-is-hidden":""},{title:_(()=>[v(s(T),{class:"mb-3",icon:"warning",color:"var(--black-500)","secondary-color":"var(--yellow-300)",size:"64"}),n(),re]),message:_(()=>[t("p",null,[n(`
        Please double check to make sure it is up and running `),s(e).baseUrl?(o(),u("span",le,[n(", and it is reachable at "),t("code",null,A(s(e).baseUrl),1)])):h("",!0)])]),_:1}))}}),ue=""+new URL("kuma-loader-v1-2aaed7d4.gif",import.meta.url).href,de=r=>(O("data-v-06e19708"),r=r(),E(),r),pe={class:"full-screen"},me={class:"loading-container"},fe=de(()=>t("img",{src:ue},null,-1)),he={class:"progress"},ge=f({__name:"AppLoadingBar",setup(r){let e;const i=$(10);return B(function(){e=window.setInterval(()=>{i.value>=100&&(window.clearInterval(e),i.value=100),i.value=Math.min(i.value+Math.ceil(Math.random()*30),100)},150)}),V(function(){window.clearInterval(e)}),(c,d)=>(o(),u("div",pe,[t("div",me,[fe,n(),t("div",he,[t("div",{style:q({width:`${i.value}%`}),class:"progress-bar",role:"progressbar","data-testid":"app-progress-bar"},null,4)])])]))}});const ve=U(ge,[["__scopeId","data-v-06e19708"]]),ye={key:0,class:"onboarding-check"},be={class:"alert-content"},Ae=f({__name:"AppOnboardingNotification",setup(r){const e=$(!1);function i(){e.value=!0}return(c,d)=>e.value===!1?(o(),u("div",ye,[v(s(W),{appearance:"success",class:"dismissible","dismiss-type":"icon",onClosed:i},{alertMessage:_(()=>[t("div",be,[t("div",null,[t("strong",null,"Welcome to "+A(s(j))+"!",1),n(` We've detected that you don't have any data plane proxies running yet. We've created an onboarding process to help you!
          `)]),n(),t("div",null,[v(s(N),{appearance:"primary",size:"small",class:"action-button",to:{name:"onboarding-welcome"}},{default:_(()=>[n(`
              Get started
            `)]),_:1})])])]),_:1})])):h("",!0)}});const Me=U(Ae,[["__scopeId","data-v-34df3ed0"]]),$e={class:"py-4"},ke=t("p",{class:"mb-4"},`
      A traffic log policy lets you collect access logs for every data plane proxy in your service mesh.
    `,-1),Se={class:"list-disc pl-4"},Ue=["href"],we=f({__name:"LoggingNotification",setup(r){const e=S();return(i,c)=>(o(),u("div",$e,[ke,n(),t("ul",Se,[t("li",null,[t("a",{href:`${s(e)("KUMA_DOCS_URL")}/policies/traffic-log/?${s(e)("KUMA_UTM_QUERY_PARAMS")}`,target:"_blank"},`
          Traffic Log policy documentation
        `,8,Ue)])])]))}}),xe={class:"py-4"},Te=t("p",{class:"mb-4"},`
      A traffic metrics policy lets you collect key data for observability of your service mesh.
    `,-1),Ne={class:"list-disc pl-4"},Ce=["href"],Le=f({__name:"MetricsNotification",setup(r){const e=S();return(i,c)=>(o(),u("div",xe,[Te,n(),t("ul",Ne,[t("li",null,[t("a",{href:`${s(e)("KUMA_DOCS_URL")}/policies/traffic-metrics/?${s(e)("KUMA_UTM_QUERY_PARAMS")}`,target:"_blank"},`
          Traffic Metrics policy documentation
        `,8,Ce)])])]))}}),Ie={class:"py-4"},Pe=t("p",{class:"mb-4"},`
      Mutual TLS (mTLS) for communication between all the components
      of your service mesh (services, control plane, data plane proxies), proxy authentication,
      and access control rules in Traffic Permissions policies all contribute to securing your mesh.
    `,-1),Re={class:"list-disc pl-4"},Ke=["href"],Be=["href"],Oe=["href"],Ee=f({__name:"MtlsNotification",setup(r){const e=S();return(i,c)=>(o(),u("div",Ie,[Pe,n(),t("ul",Re,[t("li",null,[t("a",{href:`${s(e)("KUMA_DOCS_URL")}/security/certificates/?${s(e)("KUMA_UTM_QUERY_PARAMS")}`,target:"_blank"},`
          Secure access across services
        `,8,Ke)]),n(),t("li",null,[t("a",{href:`${s(e)("KUMA_DOCS_URL")}/policies/mutual-tls/?${s(e)("KUMA_UTM_QUERY_PARAMS")}`,target:"_blank"},`
          Mutual TLS
        `,8,Be)]),n(),t("li",null,[t("a",{href:`${s(e)("KUMA_DOCS_URL")}/policies/traffic-permissions/?${s(e)("KUMA_UTM_QUERY_PARAMS")}`,target:"_blank"},`
          Traffic Permissions policy documentation
        `,8,Oe)])])]))}}),ze={class:"py-4"},De=t("p",{class:"mb-4"},`
      A traffic trace policy lets you enable tracing logs and a third-party tracing solution to send them to.
    `,-1),We={class:"list-disc pl-4"},Ye=["href"],Qe=f({__name:"TracingNotification",setup(r){const e=S();return(i,c)=>(o(),u("div",ze,[De,n(),t("ul",We,[t("li",null,[t("a",{href:`${s(e)("KUMA_DOCS_URL")}/policies/traffic-trace/?${s(e)("KUMA_UTM_QUERY_PARAMS")}`,target:"_blank"},`
          Traffic Trace policy documentation
        `,8,Ye)])])]))}}),Ve={class:"flex items-center"},qe=f({__name:"SingleMeshNotifications",setup(r){const e=k(),i={LoggingNotification:we,MetricsNotification:Le,MtlsNotification:Ee,TracingNotification:Qe};return(c,d)=>(o(),l(ne,{"multiple-open":""},{default:_(()=>[(o(!0),u(z,null,F(s(e).getters["notifications/singleMeshNotificationItems"],p=>(o(),l(oe,{key:p.name},{"accordion-header":_(()=>[t("div",Ve,[p.isCompleted?(o(),l(s(T),{key:0,color:"var(--green-500)",icon:"check",size:"20",class:"mr-4"})):(o(),l(s(T),{key:1,icon:"warning",color:"var(--black-500)","secondary-color":"var(--yellow-300)",size:"20",class:"mr-4"})),n(),t("strong",null,A(p.name),1)])]),"accordion-content":_(()=>[p.component?(o(),l(D(i[p.component]),{key:0})):(o(),l(s(ee),{key:1},{body:_(()=>[n(A(p.content),1)]),_:2},1024))]),_:2},1024))),128))]),_:1}))}}),je=r=>(O("data-v-ce28c0f7"),r=r(),E(),r),Fe={class:"mr-4"},He=je(()=>t("span",{class:"mr-2"},[t("strong",null,"Pro tip:"),n(`

            You might want to adjust your mesh configuration
          `)],-1)),Ze={key:0},Ge={class:"text-xl tracking-wide"},Je={key:1},Xe={class:"text-xl tracking-wide"},et=f({__name:"NotificationManager",setup(r){const e=k(),i=$(!0),c=y(()=>e.state.selectedMesh?e.getters["notifications/meshNotificationItemMapWithAction"][e.state.selectedMesh]:!1);B(function(){const a=R.get("hideCheckMeshAlert");i.value=a!=="yes"});function d(){i.value=!1,R.set("hideCheckMeshAlert","yes")}function p(){e.dispatch("notifications/openModal")}function M(){e.dispatch("notifications/closeModal")}return(a,m)=>(o(),u("div",null,[i.value?(o(),l(s(W),{key:0,class:"mb-4",appearance:"info","dismiss-type":"icon","data-testid":"notification-info",onClosed:d},{alertMessage:_(()=>[t("div",Fe,[He,n(),v(s(N),{appearance:"outline","data-testid":"open-modal-button",onClick:p},{default:_(()=>[n(`
            Check your mesh!
          `)]),_:1})])]),_:1})):h("",!0),n(),v(s(te),{class:"modal","is-visible":s(e).state.notifications.isOpen,title:"Notifications","text-align":"left","data-testid":"notification-modal"},{"header-content":_(()=>[t("div",null,[t("div",null,[c.value?(o(),u("span",Ze,[n(`
              Some of these features are not enabled for `),t("span",Ge,'"'+A(s(e).state.selectedMesh)+'"',1),n(` mesh. Consider implementing them.
            `)])):(o(),u("span",Je,[n(`
              Looks like `),t("span",Xe,'"'+A(s(e).state.selectedMesh)+'"',1),n(` isn't missing any features. Well done!
            `)]))])])]),"body-content":_(()=>[v(qe)]),"footer-content":_(()=>[v(s(N),{appearance:"outline","data-testid":"close-modal-button",onClick:M},{default:_(()=>[n(`
          Close
        `)]),_:1})]),_:1},8,["is-visible"])]))}});const tt=U(et,[["__scopeId","data-v-ce28c0f7"]]),st={key:1},ot={class:"app-main-content"},nt=f({__name:"App",setup(r){const[e,i]=[ae(),ie()],c=k(),d=K(),p=$(c.state.globalLoading),M=y(()=>d.path),a=y(()=>d.meta.isWizard===!0),m=y(()=>c.getters.shouldShowAppError),b=y(()=>c.getters.shouldShowNotificationManager),C=y(()=>c.getters.shouldShowOnboardingNotification),L=y(()=>c.getters.shouldShowBreadcrumbs);x(()=>c.state.globalLoading,function(g){p.value=g}),x(()=>d.meta.title,function(g){I(g)}),x(()=>c.state.pageTitle,function(g){I(g)});function I(g){const w="Kuma Manager";document.title=g?`${g} | ${w}`:w}return(g,w)=>{const P=H("router-view");return p.value||s(d).name===void 0?(o(),l(ve,{key:0})):(o(),u(z,{key:1},[a.value?h("",!0):(o(),l(s(i),{key:0})),n(),s(d).meta.onboardingProcess?(o(),u("div",st,[v(P)])):(o(),u("div",{key:2,class:Z(["app-content-container",{"is-wizard":a.value}])},[a.value?h("",!0):(o(),l(s(e),{key:0})),n(),t("main",ot,[m.value?(o(),l(_e,{key:0,"data-testid":"app-error"})):h("",!0),n(),!a.value&&b.value?(o(),l(tt,{key:1})):h("",!0),n(),!a.value&&C.value?(o(),l(Me,{key:2})):h("",!0),n(),!a.value&&L.value?(o(),l(ce,{key:3})):h("",!0),n(),(o(),l(P,{key:M.value},{default:_(({Component:Y})=>[v(G,{mode:"out-in",name:"fade"},{default:_(()=>[(o(),u("div",{key:s(d).name,class:"transition-root"},[(o(),l(D(Y)))]))]),_:2},1024)]),_:1}))])],2))],64))}}});const mt=U(nt,[["__scopeId","data-v-c0ebf95c"]]);export{mt as default};
