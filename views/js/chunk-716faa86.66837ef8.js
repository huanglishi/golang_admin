(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-716faa86"],{"0164":function(a,e,t){"use strict";t("eb53")},2188:function(a,e,t){"use strict";t.r(e);var i=function(){var a=this,e=a.$createElement,t=a._self._c||e;return t("div",{staticClass:"layerbox"},[t("a-spin",{attrs:{spinning:a.confirmLoading}},[t("a-form",{attrs:{form:a.form},on:{submit:a.handleSubmit}},[t("a-form-item",{staticClass:"roiders",attrs:{labelCol:a.labelCol,wrapperCol:a.wrapperCol,label:"分类名称",hasFeedback:""}},[t("a-input",{directives:[{name:"decorator",rawName:"v-decorator",value:["name",{initialValue:a.savedata.name,rules:[{required:!0,min:1,message:"请输入至少1个字符的名称！"}]}],expression:"['name', {initialValue:savedata.name,rules: [{required: true, min: 1, message: '请输入至少1个字符的名称！'}]}]"}],attrs:{placeholder:"起一个名字"}})],1),t("a-form-item",{staticClass:"roiders",attrs:{label:"状态:",labelCol:a.labelCol,wrapperCol:a.wrapperCol}},[t("a-radio-group",{directives:[{name:"decorator",rawName:"v-decorator",value:["status",{initialValue:a.savedata.status}],expression:"['status',{initialValue:savedata.status}]"}]},[t("a-radio",{attrs:{value:0}},[a._v("正常")]),t("a-radio",{attrs:{value:1}},[a._v("隐藏")])],1)],1),t("a-form-item",{staticStyle:{"text-align":"center","margin-top":"12px"},attrs:{wrapperCol:{span:24}}},[t("a-button",{attrs:{htmlType:"submit",type:"primary"}},[a._v(a._s(a.$t("form.basic-form.form.submit")))]),t("a-button",{staticStyle:{"margin-left":"8px"},on:{click:a.onReset}},[a._v(a._s(a.$t("form.basic-form.form.reset")))])],1)],1)],1)],1)},r=[],s=(t("b0c0"),t("d3b7"),t("4de4"),t("8ded")),l=t.n(s),n={name:"AdddevCate",components:{},data:function(){return{labelCol:{xs:{span:24},sm:{span:5}},wrapperCol:{xs:{span:24},sm:{span:16}},confirmLoading:!1,form:this.$form.createForm(this),list:[],savedata:{pid:"",name:"",status:0},id:0}},props:{layerid:{type:String,default:""},lydata:{type:Object,default:function(){return{}}},lyoption:{type:Object,default:function(){return{}}}},created:function(){var a=this;if(this.lydata.data){var e=this.lydata.data;this.$nextTick((function(){a.id=e.id,a.item_data=e,a.form.setFieldsValue({name:e.name,status:e.status})}))}else this.$nextTick((function(){a.form.setFieldsValue({status:0})}));this.all_layerid=l.a.get("layerid"),this.all_layerid||(this.all_layerid=[]),this.all_layerid.push(this.layerid),l.a.set("layerid",this.all_layerid,864e5)},methods:{onReset:function(){this.form.resetFields(),this.form.setFieldsValue({status:0})},onExpand:function(a){this.expandedKeys=a,this.autoExpandParent=!1},concatarr:function(a,e){for(var t=0;t<e.length;t++)-1==a.indexOf(e[t])&&a.push(e[t]);return a},handleSubmit:function(a){var e=this;a.preventDefault(),this.form.validateFields((function(a,t){a||(e.confirmLoading=!0,t["id"]=e.id,e.$request({url:"developer/module/add",method:"post",data:t}).then((function(a){0==a.code?(e.$message.success(a.msg),e.$layer.close(e.layerid),e.all_layerid&&e.all_layerid.length>0&&(e.all_layerid=e.all_layerid.filter((function(a){return a!=e.layerid})),l.a.set("layerid",e.all_layerid,864e5)),e.$parent.efleshCate()):e.$message.error(a.msg)})).catch((function(a){})).finally((function(){e.confirmLoading=!1})))}))}}},o=n,d=(t("0164"),t("2877")),u=Object(d["a"])(o,i,r,!1,null,"c0fb8972",null);e["default"]=u.exports},eb53:function(a,e,t){}}]);