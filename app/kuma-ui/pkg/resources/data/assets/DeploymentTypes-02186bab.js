import{s as r}from"./kongponents.es-c741eab8.js";import{O as g,a as b,b as v}from"./OnboardingPage-674cec80.js";import{c as f,d as y}from"./index-93c6e791.js";import{d as h,r as V,c as z,k as x,b as i,w as e,o as u,g as n,h as t,i as l,t as M,e as d,P as S,s as D}from"./index-0fbacd76.js";import{u as G}from"./store-5ee7e2bf.js";import{_ as C}from"./_plugin-vue_export-helper-c27b6911.js";import"./index-bc240554.js";import"./datadogLogEvents-302eea7b.js";import"./DoughnutChart-46244ceb.js";const N={class:"graph-list mb-6"},O={class:"radio-button-group"},T=h({__name:"DeploymentTypes",setup(w){const p=f(),m={standalone:y(),"multi-zone":p},c=G(),o=V("standalone"),_=z(()=>m[o.value]);return x(function(){o.value=c.getters["config/getMulticlusterStatus"]?"multi-zone":"standalone"}),(k,a)=>(u(),i(g,{"with-image":""},{header:e(()=>[n(b,null,{title:e(()=>[t(`
          Learn about deployments
        `)]),description:e(()=>[l("p",null,M(d(S))+" can be deployed in standalone or multi-zone mode.",1)]),_:1})]),content:e(()=>[l("div",N,[(u(),i(D(_.value)))]),t(),l("div",O,[n(d(r),{modelValue:o.value,"onUpdate:modelValue":a[0]||(a[0]=s=>o.value=s),name:"mode","selected-value":"standalone","data-testid":"onboarding-standalone-radio-button"},{default:e(()=>[t(`
          Standalone deployment
        `)]),_:1},8,["modelValue"]),t(),n(d(r),{modelValue:o.value,"onUpdate:modelValue":a[1]||(a[1]=s=>o.value=s),name:"mode","selected-value":"multi-zone","data-testid":"onboarding-multi-zone-radio-button"},{default:e(()=>[t(`
          Multi-zone deployment
        `)]),_:1},8,["modelValue"])])]),navigation:e(()=>[n(v,{"next-step":"onboarding-configuration-types","previous-step":"onboarding-welcome"})]),_:1}))}});const q=C(T,[["__scopeId","data-v-711049bb"]]);export{q as default};
