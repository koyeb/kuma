import{d as u,o as a,f as o,t as r,b as n,g as t,e as C,cp as b,ct as f,cu as D,cv as O,c as p,w as i,i as h,a as N,cw as A,F as E,r as _}from"./index.d0c2efa7.js";const I=u({__name:"WarningDefault",props:{payload:{type:[String,Object],required:!0}},setup(e){return(c,s)=>(a(),o("span",null,r(e.payload),1))}}),W=u({__name:"WarningEnvoyIncompatible",props:{payload:{type:Object,required:!0}},setup(e){return(c,s)=>(a(),o("span",null,[n(" Envoy ("),t("strong",null,r(e.payload.envoy),1),n(") is unsupported by the current version of Kuma DP ("),t("strong",null,r(e.payload.kumaDp),1),n(") [Requirements: "),t("strong",null,r(e.payload.requirements),1),n("]. ")]))}}),v=u({__name:"WarningZoneAndKumaDPVersionsIncompatible",props:{payload:{type:Object,required:!0}},setup(e){return(c,s)=>(a(),o("span",null,[n(" There is a mismatch between versions of Kuma DP ("),t("strong",null,r(e.payload.kumaDp),1),n(") and the Zone CP. ")]))}}),K=u({__name:"WarningUnsupportedKumaDPVersion",props:{payload:{type:Object,required:!0}},setup(e){return(c,s)=>(a(),o("span",null,[n(" Unsupported version of Kuma DP ("),t("strong",null,r(e.payload.kumaDp),1),n(") ")]))}}),V=u({__name:"WarningZoneAndGlobalCPSVersionsIncompatible",props:{payload:{type:Object,required:!0}},setup(e){return(c,s)=>(a(),o("span",null,[n(" There is mismatch between versions of Zone CP ("),t("strong",null,r(e.payload.zoneCpVersion),1),n(") and the Global CP ("),t("strong",null,r(e.payload.globalCpVersion),1),n(") ")]))}}),x={name:"WarningsWidget",props:{warnings:{type:Array,required:!0}},methods:{getWarningComponent(e=""){switch(e){case O:return W;case D:return K;case f:return v;case b:return V;default:return I}}}};function B(e,c,s,S,w,d){const m=_("KAlert"),y=_("KCard");return a(),p(y,{"border-variant":"noBorder"},{body:i(()=>[t("ul",null,[(a(!0),o(E,null,h(s.warnings,({kind:l,payload:g,index:P})=>(a(),o("li",{key:`${l}/${P}`,class:"mb-1"},[N(m,{appearance:"warning"},{alertMessage:i(()=>[(a(),p(A(d.getWarningComponent(l)),{payload:g},null,8,["payload"]))]),_:2},1024)]))),128))])]),_:1})}const $=C(x,[["render",B]]);export{$ as W};
