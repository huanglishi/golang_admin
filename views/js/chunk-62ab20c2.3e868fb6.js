(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-62ab20c2"],{"1da3":function(e,a,t){"use strict";t.r(a);var r=function(){var e=this,a=e.$createElement,t=e._self._c||a;return t("div",{staticClass:"layerbox"},[t("a-spin",{attrs:{spinning:e.confirmLoading}},[t("a-form",{attrs:{form:e.form},on:{submit:e.handleSubmit}},[t("a-form-item",{staticClass:"roiders",attrs:{"label-col":e.labelCol,"wrapper-col":e.wrapperCol,label:"选择父级","has-feedback":""}},[t("a-select",{directives:[{name:"decorator",rawName:"v-decorator",value:["pid",{rules:[{required:!0,message:"请选择父级！"}]}],expression:"['pid', {rules: [{required: true, message: '请选择父级！'}]}]"}],attrs:{allowClear:"",placeholder:"选择上一级"}},[t("a-select-option",{key:0,attrs:{value:0}},[e._v("无")]),e._l(e.list,(function(a){return t("a-select-option",{key:a.id,attrs:{value:a.id}},[t("span",{domProps:{innerHTML:e._s(a.name_txt)}})])}))],2)],1),t("a-form-item",{staticClass:"roiders",attrs:{labelCol:e.labelCol,wrapperCol:e.wrapperCol,label:"分类名称",hasFeedback:""}},[t("a-input",{directives:[{name:"decorator",rawName:"v-decorator",value:["name",{initialValue:e.savedata.name,rules:[{required:!0,min:1,message:"请输入至少1个字符的名称！"}]}],expression:"['name', {initialValue:savedata.name,rules: [{required: true, min: 1, message: '请输入至少1个字符的名称！'}]}]"}],attrs:{placeholder:"起一个名字"}})],1),t("a-form-item",{staticClass:"roiders",attrs:{label:"状态:",labelCol:e.labelCol,wrapperCol:e.wrapperCol}},[t("a-radio-group",{directives:[{name:"decorator",rawName:"v-decorator",value:["status",{initialValue:e.savedata.status}],expression:"['status',{initialValue:savedata.status}]"}]},[t("a-radio",{attrs:{value:0}},[e._v("正常")]),t("a-radio",{attrs:{value:1}},[e._v("隐藏")])],1)],1),t("a-form-item",{staticStyle:{"text-align":"center","margin-top":"12px"},attrs:{wrapperCol:{span:24}}},[t("a-button",{attrs:{htmlType:"submit",type:"primary"}},[e._v(e._s(e.$t("form.basic-form.form.submit")))]),t("a-button",{staticStyle:{"margin-left":"8px"},on:{click:e.onReset}},[e._v(e._s(e.$t("form.basic-form.form.reset")))])],1)],1)],1)],1)},i=[],s=(t("d3b7"),t("b0c0"),t("4de4"),t("8ded")),l=t.n(s),n={name:"AdddevCate",components:{},data:function(){return{labelCol:{xs:{span:24},sm:{span:5}},wrapperCol:{xs:{span:24},sm:{span:16}},confirmLoading:!1,form:this.$form.createForm(this),list:[],savedata:{pid:"",name:"",status:0},id:0}},props:{layerid:{type:String,default:""},lydata:{type:Object,default:function(){return{}}},lyoption:{type:Object,default:function(){return{}}}},created:function(){var e=this;if(this.$request({url:"developer/apidoc/getparentcate",method:"get"}).then((function(a){e.list=a.data})).catch((function(e){})).finally((function(){})),this.lydata.data){var a=this.lydata.data;this.$nextTick((function(){e.id=a.id,e.item_data=a,e.form.setFieldsValue({pid:a.pid,name:a.name,status:a.status})}))}else this.$nextTick((function(){e.form.setFieldsValue({status:0})}));this.all_layerid=l.a.get("layerid"),this.all_layerid||(this.all_layerid=[]),this.all_layerid.push(this.layerid),l.a.set("layerid",this.all_layerid,864e5)},methods:{onReset:function(){this.form.resetFields(),this.form.setFieldsValue({status:0})},onExpand:function(e){this.expandedKeys=e,this.autoExpandParent=!1},concatarr:function(e,a){for(var t=0;t<a.length;t++)-1==e.indexOf(a[t])&&e.push(a[t]);return e},handleSubmit:function(e){var a=this;e.preventDefault(),this.form.validateFields((function(e,t){e||(a.confirmLoading=!0,t["id"]=a.id,a.$request({url:"developer/apidoc/add",method:"post",data:t}).then((function(e){0==e.code?(a.$message.success(e.msg),a.$layer.close(a.layerid),a.all_layerid&&a.all_layerid.length>0&&(a.all_layerid=a.all_layerid.filter((function(e){return e!=a.layerid})),l.a.set("layerid",a.all_layerid,864e5)),a.$parent.efleshCate()):a.$message.error(e.msg)})).catch((function(e){})).finally((function(){a.confirmLoading=!1})))}))}}},o=n,d=(t("313b"),t("2877")),u=Object(d["a"])(o,r,i,!1,null,"2792efae",null);a["default"]=u.exports},"313b":function(e,a,t){"use strict";t("3c50")},"3c50":function(e,a,t){}}]);