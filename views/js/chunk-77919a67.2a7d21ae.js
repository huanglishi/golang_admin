(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-77919a67","chunk-7373300a","chunk-2d0b6517"],{"00d8":function(e,t){(function(){var t="ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/",a={rotl:function(e,t){return e<<t|e>>>32-t},rotr:function(e,t){return e<<32-t|e>>>t},endian:function(e){if(e.constructor==Number)return 16711935&a.rotl(e,8)|4278255360&a.rotl(e,24);for(var t=0;t<e.length;t++)e[t]=a.endian(e[t]);return e},randomBytes:function(e){for(var t=[];e>0;e--)t.push(Math.floor(256*Math.random()));return t},bytesToWords:function(e){for(var t=[],a=0,r=0;a<e.length;a++,r+=8)t[r>>>5]|=e[a]<<24-r%32;return t},wordsToBytes:function(e){for(var t=[],a=0;a<32*e.length;a+=8)t.push(e[a>>>5]>>>24-a%32&255);return t},bytesToHex:function(e){for(var t=[],a=0;a<e.length;a++)t.push((e[a]>>>4).toString(16)),t.push((15&e[a]).toString(16));return t.join("")},hexToBytes:function(e){for(var t=[],a=0;a<e.length;a+=2)t.push(parseInt(e.substr(a,2),16));return t},bytesToBase64:function(e){for(var a=[],r=0;r<e.length;r+=3)for(var n=e[r]<<16|e[r+1]<<8|e[r+2],o=0;o<4;o++)8*r+6*o<=8*e.length?a.push(t.charAt(n>>>6*(3-o)&63)):a.push("=");return a.join("")},base64ToBytes:function(e){e=e.replace(/[^A-Z0-9+\/]/gi,"");for(var a=[],r=0,n=0;r<e.length;n=++r%4)0!=n&&a.push((t.indexOf(e.charAt(r-1))&Math.pow(2,-2*n+8)-1)<<2*n|t.indexOf(e.charAt(r))>>>6-2*n);return a}};e.exports=a})()},"044b":function(e,t){function a(e){return!!e.constructor&&"function"===typeof e.constructor.isBuffer&&e.constructor.isBuffer(e)}function r(e){return"function"===typeof e.readFloatLE&&"function"===typeof e.slice&&a(e.slice(0,0))}
/*!
 * Determine if an object is a Buffer
 *
 * @author   Feross Aboukhadijeh <https://feross.org>
 * @license  MIT
 */
e.exports=function(e){return null!=e&&(a(e)||r(e)||!!e._isBuffer)}},"1d3a":function(e,t,a){"use strict";a.r(t),a.d(t,"getdata",(function(){return n})),a.d(t,"getdatelist",(function(){return o})),a.d(t,"addscancode",(function(){return i}));var r=a("b775");function n(e){return Object(r["b"])({url:"admin/user/list",method:"get",params:e})}function o(e){return Object(r["b"])({url:"/onepage.purchasefile/getdatelist",method:"get",params:e})}function i(e){return Object(r["b"])({url:"/onepage.purchasefile/addscancode",method:"post",data:e})}},"209e":function(e,t,a){"use strict";a.r(t);var r=function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("page-header-wrapper",{attrs:{title:" "}},[a("div",{staticClass:"table-operator"},[a("a-button",{directives:[{name:"action",rawName:"v-action:add",arg:"add"}],attrs:{type:"primary",icon:"plus"},on:{click:function(t){return e.onAddForm()}}},[e._v("新建")]),a("a-dropdown",{attrs:{disabled:0==e.selectedRowKeys.length}},[a("a-menu",{attrs:{slot:"overlay"},slot:"overlay"},[a("a-menu-item",{directives:[{name:"action",rawName:"v-action:delete",arg:"delete"}],key:"1",on:{click:function(t){return e.del_rule(e.selectedRowKeys)}}},[a("a-icon",{attrs:{type:"delete"}}),e._v("删除")],1),a("a-menu-item",{directives:[{name:"action",rawName:"v-action:update",arg:"update"}],key:"2",on:{click:function(t){return e.onlock(2)}}},[a("a-icon",{attrs:{type:"lock"}}),e._v("锁定")],1),a("a-menu-item",{directives:[{name:"action",rawName:"v-action:update",arg:"update"}],key:"3",on:{click:function(t){return e.onlock(0)}}},[a("a-icon",{attrs:{type:"unlock"}}),e._v("启用")],1)],1),a("a-button",{staticStyle:{"margin-left":"8px"}},[e._v(" 批量操作 "),a("a-icon",{attrs:{type:"down"}})],1)],1)],1),a("div",{staticClass:"tablebox"},[a("s-table",{ref:"table",attrs:{bordered:"",size:"default",rowKey:"id",columns:e.columns,data:e.loadData,rowSelection:e.rowSelection,showPagination:"auto"},scopedSlots:e._u([{key:"groupname",fn:function(t){return e._l(t,(function(t){return a("a-tag",{key:t.id,staticStyle:{"user-select":"none"},attrs:{color:"orange"}},[e._v(e._s(t.name))])}))}},{key:"DateTime",fn:function(t){return a("span",{},[e._v(" "+e._s(e._f("timestampToTime")(t))+" ")])}},{key:"status",fn:function(t){return a("span",{},[a("a-tag",{staticStyle:{"user-select":"none"},attrs:{color:e._f("statusTypeFilter")(t)}},[e._v(e._s(e._f("statusFilter")(t)))])],1)}},{key:"action",fn:function(t,r){return a("span",{},[[a("a",{directives:[{name:"action",rawName:"v-action:edit",arg:"edit"}],on:{click:function(t){return e.onAddForm(r)}}},[e._v("编辑")]),a("a-divider",{directives:[{name:"action",rawName:"v-action:edit",arg:"edit"},{name:"action",rawName:"v-action:delete",arg:"delete"}],attrs:{type:"vertical"}}),a("a-popconfirm",{directives:[{name:"action",rawName:"v-action:delete",arg:"delete"}],attrs:{title:e.$t("app.table.del.aq")},on:{confirm:function(){return e.del_rule([r.id])}}},[a("a",{attrs:{href:"javascript:;"}},[e._v(e._s(e.$t("app.table.del")))])])]],2)}}])})],1)])},n=[],o=(a("4de4"),a("d3b7"),a("2af9")),i=a("4e38"),s=a("1d3a"),l=a("8ded"),u=a.n(l),c=[{title:"ID",maxwidth:"60px",dataIndex:"id"},{title:"姓名",dataIndex:"name"},{title:"用户账号",dataIndex:"username"},{title:"手机号码",dataIndex:"telephone"},{title:"权限组",dataIndex:"groupname",scopedSlots:{customRender:"groupname"}},{title:"状态",dataIndex:"status",scopedSlots:{customRender:"status"}},{title:"最后登录IP",dataIndex:"lastLoginIp"},{title:"最后登录时间",dataIndex:"lastLoginTime",scopedSlots:{customRender:"DateTime"},sorter:!0},{title:"操作",dataIndex:"action",width:"150px",scopedSlots:{customRender:"action"}}],d={0:{status:"",text:"未登录"},1:{status:"blue",text:"登录中"},2:{status:"pink",text:"锁定"}},f={components:{STable:o["m"]},data:function(){var e=this;return this.columns=c,{queryParam:{},loadData:function(t){var a=Object.assign({},t,e.queryParam);return Object(s["getdata"])(a).then((function(t){if(0==t.code)return t.data;e.$message.error(t.msg)}))},selectedRowKeys:[],selectedRows:[]}},computed:{rowSelection:function(){return{selectedRowKeys:this.selectedRowKeys,onChange:this.onSelectChange}}},methods:{onAddForm:function(e){this.layerid=this.$layer.iframe({content:{content:i["default"],parent:this,data:{data:e||null}},area:["800px","600px"],title:"添加/编辑菜单",maxmin:!0,shade:!1,shadeClose:!1,cancel:function(e){var t=u.a.get("layerid");t&&t.length>0&&(t=t.filter((function(t){return t!=e})),u.a.set("layerid",t,864e5))}})},reflesh:function(){this.$refs.table.refresh(!0)},onlock:function(e){var t=this,a=this.$message.loading("更新中，请稍后...",0);this.$request({url:"admin/user/updata",method:"post",data:{status:e,ids:this.selectedRowKeys}}).then((function(e){0==e.code?(t.$message.success(e.msg),t.$refs.table.refresh(!0)):t.$message.error(e.msg)})).catch((function(e){})).finally((function(){setTimeout(a,1)}))},del_rule:function(e){var t=this,a=this.$message.loading("删除中，请稍后...",0);this.$request({url:"admin/user/del",method:"delete",data:{ids:e}}).then((function(e){0==e.code?(t.$message.success(e.msg),t.$refs.table.refresh(!0)):t.$message.error(e.msg)})).catch((function(e){t.$message.error("网络错误，请稍后再试")})).finally((function(){setTimeout(a,1)}))},onSelectChange:function(e,t){this.selectedRowKeys=e,this.selectedRows=t}},filters:{timestampToTime:function(e){if(!e)return"---";var t=new Date(1e3*e),a=t.getFullYear()+"-",r=t.getMonth()+1+"-",n=t.getDate()+" ",o=t.getHours()+":",i=t.getMinutes();return a+r+n+o+i},statusFilter:function(e){return d[e].text},statusTypeFilter:function(e){return d[e].status}},beforeRouteLeave:function(e,t,a){var r=u.a.get("layerid");if(r&&r.length>0)for(var n in r)try{this.$layer.min(r[n])}catch(o){continue}a()}},m=f,p=(a("ce3f"),a("2877")),h=Object(p["a"])(m,r,n,!1,null,"d2b17a60",null);t["default"]=h.exports},"4e38":function(e,t,a){"use strict";a.r(t);var r=function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("div",{staticClass:"layerbox"},[a("a-spin",{attrs:{spinning:e.confirmLoading}},[a("a-form",{attrs:{form:e.form},on:{submit:e.handleSubmit}},[a("a-form-item",{attrs:{label:e.$t("form.basic-form.name"),labelCol:e.labelCol,wrapperCol:e.wrapperCol,hasFeedback:""}},[a("a-input",{directives:[{name:"decorator",rawName:"v-decorator",value:["name",{rules:[{required:!0,message:e.$t("form.basic-form.name.placeholder")}]}],expression:"['name', {rules:[{required: true, message: $t('form.basic-form.name.placeholder')}]}]"}],attrs:{allowClear:"",placeholder:e.$t("form.basic-form.name.placeholder")}})],1),a("a-form-item",{attrs:{label:e.$t("form.basic-form.phone-number"),labelCol:e.labelCol,wrapperCol:e.wrapperCol,hasFeedback:""}},[a("a-input",{directives:[{name:"decorator",rawName:"v-decorator",value:["telephone",{rules:[{pattern:/^1([358][0-9]|4[579]|66|7[0135678]|9[89])[0-9]{8}$/gi,message:e.$t("form.basic-form.phone-number.required")}]}],expression:"['telephone',{rules: [\n          {pattern: /^1([358][0-9]|4[579]|66|7[0135678]|9[89])[0-9]{8}$/ig, message:$t('form.basic-form.phone-number.required')}]}]"}],attrs:{allowClear:"",placeholder:e.$t("form.basic-form.phone-number.required")}})],1),a("a-form-item",{attrs:{label:e.$t("user.login.account"),labelCol:e.labelCol,wrapperCol:e.wrapperCol,hasFeedback:""}},[a("a-input",{directives:[{name:"decorator",rawName:"v-decorator",value:["username",{rules:[{required:!0,message:"请输入登录账户",pattern:/^[A-Za-z0-9_\-_\/_@]+$/gi,message:"输入只能是数字、字母、下划线"}]}],expression:"[\n                'username',\n                 {rules: [{ required: true, message:'请输入登录账户' ,pattern: /^[A-Za-z0-9_\\-_\\/_@]+$/ig, message: '输入只能是数字、字母、下划线'}]}]"}],attrs:{allowClear:"",type:"text",placeholder:e.$t("user.login.username.placeholder")}},[a("a-icon",{style:{color:"rgba(0,0,0,.25)"},attrs:{slot:"prefix",type:"user"},slot:"prefix"})],1)],1),a("a-form-item",{attrs:{label:e.$t("user.login.password"),labelCol:e.labelCol,wrapperCol:e.wrapperCol,help:e.id>0?"不修改密码请留空":"不填写密码则默认123456",hasFeedback:""}},[a("a-input-password",{directives:[{name:"decorator",rawName:"v-decorator",value:["password",{validateTrigger:"blur"}],expression:"[\n                'password',\n                { validateTrigger: 'blur'}\n              ]"}],attrs:{allowClear:"",placeholder:e.$t("user.login.password.placeholder")}},[a("a-icon",{style:{color:"rgba(0,0,0,.25)"},attrs:{slot:"prefix",type:"lock"},slot:"prefix"})],1)],1),a("a-form-item",{attrs:{label:"有效时间",labelCol:e.labelCol,wrapperCol:e.wrapperCol,help:"不选择永久有效，若选择则以选择时间为有效时限"}},[a("a-date-picker",{directives:[{name:"decorator",rawName:"v-decorator",value:["valid_time"],expression:"['valid_time']"}],staticStyle:{width:"100%"},attrs:{format:e.dateFormat}})],1),a("a-form-item",{attrs:{labelCol:e.labelCol,wrapperCol:e.wrapperCol,label:"权限组",hasFeedback:""}},[a("a-select",{directives:[{name:"decorator",rawName:"v-decorator",value:["groupids",{rules:[{required:!0,message:"请选择权限组"}]}],expression:"['groupids',{rules:[{required: true, message: '请选择权限组'}]}]"}],staticStyle:{width:"100%"},attrs:{mode:"multiple",allowClear:!0,placeholder:"选择权限组"}},e._l(e.list,(function(t,r){return a("a-select-option",{key:r,attrs:{value:t.id}},[a("span",{domProps:{innerHTML:e._s(t.name_txt)}})])})),1)],1),a("a-form-item",{staticStyle:{"text-align":"center","margin-top":"12px"},attrs:{wrapperCol:{span:24}}},[a("a-button",{attrs:{htmlType:"submit",type:"primary"}},[e._v(e._s(e.$t("form.basic-form.form.submit")))]),a("a-button",{staticStyle:{"margin-left":"8px"},on:{click:e.onReset}},[e._v(e._s(e.$t("form.basic-form.form.reset")))])],1)],1)],1)],1)},n=[],o=a("5530"),i=(a("b0c0"),a("d3b7"),a("4de4"),a("c1df")),s=a.n(i),l=a("6821"),u=a.n(l),c=a("8ded"),d=a.n(c),f={data:function(){return{labelCol:{xs:{span:24},sm:{span:5}},wrapperCol:{xs:{span:24},sm:{span:16}},confirmLoading:!1,form:this.$form.createForm(this),dateFormat:"YYYY-MM-DD",list:[],savedata:{},id:0,all_layerid:[]}},props:{layerid:{type:String,default:""},lydata:{type:Object,default:function(){return{}}},lyoption:{type:Object,default:function(){return{}}}},created:function(){var e=this;if(this.getdata(),this.lydata.data){var t=this.lydata.data;this.$nextTick((function(){e.id=t.id;var a="";0!=t.valid_time&&(a={valid_time:s()(e.dateToStr(1e3*t.valid_time),e.dateFormat)}),e.form.setFieldsValue(Object(o["a"])({name:t.name,username:t.username,telephone:t.telephone,groupids:t.groupids},a))}))}else this.$nextTick((function(){e.form.setFieldsValue({})}));this.all_layerid=d.a.get("layerid"),this.all_layerid||(this.all_layerid=[]),this.all_layerid.push(this.layerid),d.a.set("layerid",this.all_layerid,864e5)},methods:{getdata:function(){var e=this;this.$request({url:"auth/group/grouptree",method:"get"}).then((function(t){e.list=t.data})).catch((function(e){})).finally((function(){}))},onReset:function(){this.form.resetFields()},handleSubmit:function(e){var t=this;e.preventDefault(),this.form.validateFields((function(e,a){e||(t.confirmLoading=!0,a["id"]=t.id,0==t.id&&(a["password"]||(a["password"]="123456")),a["password"]&&(a["password"]=u()(a["password"])),a.valid_time&&(a["valid_time"]=t.dateToStr(a.valid_time)),t.$request({url:"admin/user/add",method:"post",data:a}).then((function(e){0==e.code?(t.$message.success(e.msg),t.$layer.close(t.layerid),t.all_layerid&&t.all_layerid.length>0&&(t.all_layerid=t.all_layerid.filter((function(e){return e!=t.layerid})),d.a.set("layerid",t.all_layerid,864e5)),t.$parent.reflesh()):t.$message.error(e.msg)})).catch((function(e){})).finally((function(){t.confirmLoading=!1})))}))},dateToStr:function(e){var t=new Date(e),a=t.getFullYear(),r=t.getMonth()+1,n=t.getDate();return r<10&&(r="0"+r),n<10&&(n="0"+n),a+"/"+r+"/"+n}}},m=f,p=(a("c6ca"),a("2877")),h=Object(p["a"])(m,r,n,!1,null,"46d817a8",null);t["default"]=h.exports},6821:function(e,t,a){(function(){var t=a("00d8"),r=a("9a634").utf8,n=a("044b"),o=a("9a634").bin,i=function(e,a){e.constructor==String?e=a&&"binary"===a.encoding?o.stringToBytes(e):r.stringToBytes(e):n(e)?e=Array.prototype.slice.call(e,0):Array.isArray(e)||e.constructor===Uint8Array||(e=e.toString());for(var s=t.bytesToWords(e),l=8*e.length,u=1732584193,c=-271733879,d=-1732584194,f=271733878,m=0;m<s.length;m++)s[m]=16711935&(s[m]<<8|s[m]>>>24)|4278255360&(s[m]<<24|s[m]>>>8);s[l>>>5]|=128<<l%32,s[14+(l+64>>>9<<4)]=l;var p=i._ff,h=i._gg,g=i._hh,y=i._ii;for(m=0;m<s.length;m+=16){var v=u,b=c,w=d,_=f;u=p(u,c,d,f,s[m+0],7,-680876936),f=p(f,u,c,d,s[m+1],12,-389564586),d=p(d,f,u,c,s[m+2],17,606105819),c=p(c,d,f,u,s[m+3],22,-1044525330),u=p(u,c,d,f,s[m+4],7,-176418897),f=p(f,u,c,d,s[m+5],12,1200080426),d=p(d,f,u,c,s[m+6],17,-1473231341),c=p(c,d,f,u,s[m+7],22,-45705983),u=p(u,c,d,f,s[m+8],7,1770035416),f=p(f,u,c,d,s[m+9],12,-1958414417),d=p(d,f,u,c,s[m+10],17,-42063),c=p(c,d,f,u,s[m+11],22,-1990404162),u=p(u,c,d,f,s[m+12],7,1804603682),f=p(f,u,c,d,s[m+13],12,-40341101),d=p(d,f,u,c,s[m+14],17,-1502002290),c=p(c,d,f,u,s[m+15],22,1236535329),u=h(u,c,d,f,s[m+1],5,-165796510),f=h(f,u,c,d,s[m+6],9,-1069501632),d=h(d,f,u,c,s[m+11],14,643717713),c=h(c,d,f,u,s[m+0],20,-373897302),u=h(u,c,d,f,s[m+5],5,-701558691),f=h(f,u,c,d,s[m+10],9,38016083),d=h(d,f,u,c,s[m+15],14,-660478335),c=h(c,d,f,u,s[m+4],20,-405537848),u=h(u,c,d,f,s[m+9],5,568446438),f=h(f,u,c,d,s[m+14],9,-1019803690),d=h(d,f,u,c,s[m+3],14,-187363961),c=h(c,d,f,u,s[m+8],20,1163531501),u=h(u,c,d,f,s[m+13],5,-1444681467),f=h(f,u,c,d,s[m+2],9,-51403784),d=h(d,f,u,c,s[m+7],14,1735328473),c=h(c,d,f,u,s[m+12],20,-1926607734),u=g(u,c,d,f,s[m+5],4,-378558),f=g(f,u,c,d,s[m+8],11,-2022574463),d=g(d,f,u,c,s[m+11],16,1839030562),c=g(c,d,f,u,s[m+14],23,-35309556),u=g(u,c,d,f,s[m+1],4,-1530992060),f=g(f,u,c,d,s[m+4],11,1272893353),d=g(d,f,u,c,s[m+7],16,-155497632),c=g(c,d,f,u,s[m+10],23,-1094730640),u=g(u,c,d,f,s[m+13],4,681279174),f=g(f,u,c,d,s[m+0],11,-358537222),d=g(d,f,u,c,s[m+3],16,-722521979),c=g(c,d,f,u,s[m+6],23,76029189),u=g(u,c,d,f,s[m+9],4,-640364487),f=g(f,u,c,d,s[m+12],11,-421815835),d=g(d,f,u,c,s[m+15],16,530742520),c=g(c,d,f,u,s[m+2],23,-995338651),u=y(u,c,d,f,s[m+0],6,-198630844),f=y(f,u,c,d,s[m+7],10,1126891415),d=y(d,f,u,c,s[m+14],15,-1416354905),c=y(c,d,f,u,s[m+5],21,-57434055),u=y(u,c,d,f,s[m+12],6,1700485571),f=y(f,u,c,d,s[m+3],10,-1894986606),d=y(d,f,u,c,s[m+10],15,-1051523),c=y(c,d,f,u,s[m+1],21,-2054922799),u=y(u,c,d,f,s[m+8],6,1873313359),f=y(f,u,c,d,s[m+15],10,-30611744),d=y(d,f,u,c,s[m+6],15,-1560198380),c=y(c,d,f,u,s[m+13],21,1309151649),u=y(u,c,d,f,s[m+4],6,-145523070),f=y(f,u,c,d,s[m+11],10,-1120210379),d=y(d,f,u,c,s[m+2],15,718787259),c=y(c,d,f,u,s[m+9],21,-343485551),u=u+v>>>0,c=c+b>>>0,d=d+w>>>0,f=f+_>>>0}return t.endian([u,c,d,f])};i._ff=function(e,t,a,r,n,o,i){var s=e+(t&a|~t&r)+(n>>>0)+i;return(s<<o|s>>>32-o)+t},i._gg=function(e,t,a,r,n,o,i){var s=e+(t&r|a&~r)+(n>>>0)+i;return(s<<o|s>>>32-o)+t},i._hh=function(e,t,a,r,n,o,i){var s=e+(t^a^r)+(n>>>0)+i;return(s<<o|s>>>32-o)+t},i._ii=function(e,t,a,r,n,o,i){var s=e+(a^(t|~r))+(n>>>0)+i;return(s<<o|s>>>32-o)+t},i._blocksize=16,i._digestsize=16,e.exports=function(e,a){if(void 0===e||null===e)throw new Error("Illegal argument "+e);var r=t.wordsToBytes(i(e,a));return a&&a.asBytes?r:a&&a.asString?o.bytesToString(r):t.bytesToHex(r)}})()},"9a634":function(e,t){var a={utf8:{stringToBytes:function(e){return a.bin.stringToBytes(unescape(encodeURIComponent(e)))},bytesToString:function(e){return decodeURIComponent(escape(a.bin.bytesToString(e)))}},bin:{stringToBytes:function(e){for(var t=[],a=0;a<e.length;a++)t.push(255&e.charCodeAt(a));return t},bytesToString:function(e){for(var t=[],a=0;a<e.length;a++)t.push(String.fromCharCode(e[a]));return t.join("")}}};e.exports=a},a0f9:function(e,t,a){},c6ca:function(e,t,a){"use strict";a("a0f9")},ce3f:function(e,t,a){"use strict";a("ea79")},ea79:function(e,t,a){}}]);