<?php
/**
 * +----------------------------------------------------------------------
 * | 通用内容控制器
 * +----------------------------------------------------------------------
 *                      .::::.
 *                    .::::::::.            | AUTHOR: siyu
 *                    :::::::::::           | EMAIL: 407593529@qq.com
 *                 ..:::::::::::'           | QQ: 407593529
 *             '::::::::::::'               | WECHAT: zhaoyingjie4125
 *                .::::::::::               | DATETIME: 2019/03/28
 *           '::::::::::::::..
 *                ..::::::::::::.
 *              ``::::::::::::::::
 *               ::::``:::::::::'        .:::.
 *              ::::'   ':::::'       .::::::::.
 *            .::::'      ::::     .:::::::'::::.
 *           .:::'       :::::  .:::::::::' ':::::.
 *          .::'        :::::.:::::::::'      ':::::.
 *         .::'         ::::::::::::::'         ``::::.
 *     ...:::           ::::::::::::'              ``::.
 *   ```` ':.          ':::::::::'                  ::::..
 *                      '.:::::'                    ':'````..
 * +----------------------------------------------------------------------
 */
namespace app\home\controller;
use app\common\model\Cate;
use app\common\model\Module as M;
use app\common\model\System;
use think\Db;
use think\facade\Request;
use think\captcha\Captcha;
use think\facade\Cache;
class Error extends Base
{
    public function initialize()
    {
        parent::initialize();
        //当前模块
        $this->module = strtolower(Request::module());
        //当前控制器
        $this->table = strtolower(Request::controller());
        //模型ID
//        $this->catedata=Cate::where('id','=',input('catId'))->field('module_id,template_list,template_show')->find();
        $this->moduleid =Cate::where('id','=',input('catId'))->value("module_id");
        //当前表
        $this->moduledata = Db::name("website_module")->where('id','=',$this->moduleid)->field('name,model_name,template_list,template_show,singlepage')->find();
        //系统配置
        $this->system = System::where('id',$this->request->param('site_id'))->find();
        if($this->system['wxqrcode']){
            $this->system['wxqrcode']= $this->system['file_root']. $this->system['wxqrcode'];
        }
        if($this->system['logo']){
            $this->system['logo']= $this->system['file_root']. $this->system['logo'];
        }
    }

    //列表
    public function index(){
        if(Request::param('catId')){
            $cate = Cate::where('id','=',Request::param('catId'))
                ->find();
            $cate['topid'] = $cate['pid']?$cate['pid']:$cate['id'];
        }else{
             $this->redirect("/");
        }

        if($this->moduledata["singlepage"]==1){
            //单页模型
            //查找第一条记录
            $info = Db::name($this->moduledata["model_name"])
                ->where('cid','=',Request::param('catId'))
                ->find();
            $this->view->assign('info', $info);//单页内容
            //定义tdk
            $title  = $cate['title'] ? $cate['title'] :
                ($cate['catname'] ? $cate['catname'] : $info['title']);
            $keywords    = $cate['keywords'] ? $cate['keywords'] : $this->system['keyword'];
            $description = $cate['description'] ? $cate['description'] : $this->system['description'];
        }else{
            //列表模型
            //定义tdk
            $title       = $cate['title']       ? $cate['title']       : $cate['catname'];
            $keywords    = $cate['keywords']    ? $cate['keywords']    : $this->system['keyword'];
            $description = $cate['description'] ? $cate['description'] : $this->system['description'];
        }

        $site_id=$this->request->param('site_id');
        $html=$this->system['html'];
        $this->view->assign('tplappend',"site_$site_id/$html/");//站点id
        $this->view->assign('cate',        $cate);//栏目信息
        $this->view->assign('system',      $this->system);//系统信息
        $this->view->assign('site_id',     $site_id);//站点id
        if(date("Y",time())>date("Y",$this->system["copyrigstart"])){
            $copyrightdate= $this->system["copyrigstart"]."-".date("Y",time());
        }else{
            $copyrightdate= date("Y",time());
        }
        $this->view->assign('copyrightdate', $copyrightdate);//版权时间
        $this->view->assign('public',      '/template/site_'.$site_id.'/');//公共目录
        $this->view->assign('title',       $title);//seo信息
        $this->view->assign('keywords',    $keywords);//seo信息
        $this->view->assign('description', $description);//seo信息
        //分页
        $catId=Request::param('catId');
        $__ALLCATE__ = Db::name('website_article_cate')->field('id,pid')->select();
        $__IDS__ = getChildsIdStr(getChildsId($__ALLCATE__,$catId),$catId);
        $totalCount = Db::name($this->moduledata["model_name"])
            ->where('cid','in',$__IDS__)
            ->count();
        $totalPage=ceil($totalCount/$cate["pagesize"]);
        $pagelist=[];
        for ($item=1; $item<=$totalPage; $item++){
            array_push($pagelist,$item);
        }
        $request = Request::instance();
        $page=Request::param('page');
        if(empty($page)){
            $page=1;
        }
        $pagedata=["baseUrl"=>$request->baseUrl(),"pagelist"=>$pagelist,"totalCount"=>$totalCount,"pagesize"=>$cate["pagesize"],
            "totalPage"=>$totalPage,"page"=>$page];
//        echo "日志：". $request->baseUrl();
        $this->view->assign('pagedata', $pagedata);//分页
        $template="";
        if($cate['template_list']){
            $template=str_replace('.html', '', $cate['template_list']);
        }elseif($this->moduledata["template_list"]){
            $template=str_replace('.html', '',$this->moduledata["template_list"]);
        }
        if(empty($template)){
            $template=$this->table.'_list';
        }
//        echo  "列表";
        $template="./template/site_$site_id/$html/$template.html";
        return $this->view->fetch($template);
    }

    //详情
    public function info(){
        if(Request::param('id')){
            //点击数增加
            Db::name($this->moduledata["model_name"])
                ->where('id',Request::param('id'))
                ->setInc('visits');
        }
        if(Request::param('catId')){
            $cate = Cate::where('id','=',Request::param('catId'))->find();
            $cate['topid'] = $cate['pid']?$cate['pid']:$cate['id'];
        }
        //查找详情信息
        $info = Db::name($this->moduledata["model_name"])
            ->where('id',input('id'))
            ->find();
//        $info = changefield($info,$this->moduleid);
        $info['url'] = getShowUrl($info);

        //定义tdk
        $title  = $cate['title'] ? $cate['title'] : $cate['catname'];
        $keywords = $info['keywords']?$info['keywords']:($cate['keywords'] ? $cate['keywords'] : $this->system['keyword']);
        $description = $info['des']?$info['des']:($cate['description'] ? $cate['description'] : $this->system['des']);

        $site_id=$this->request->param('site_id');
        $html=$this->system['html'];
        $this->view->assign('tplappend',"site_$site_id/$html/");//站点id
        $this->view->assign('cate'        , $cate);        //栏目信息
        $this->view->assign('info'        , $info);        //详情信息
        $this->view->assign('system'      , $this->system);      //系统信息
        $this->view->assign('site_id',     $site_id);//站点id
        if(date("Y",time())>date("Y",$this->system["copyrigstart"])){
            $copyrightdate= $this->system["copyrigstart"]."-".date("Y",time());
        }else{
            $copyrightdate= date("Y",time());
        }
        $this->view->assign('copyrightdate', $copyrightdate);//版权时间
        $this->view->assign('public'      , '/template/site_'.$site_id.'/');//公共目录
        $this->view->assign('title'       , $title);       //seo信息
        $this->view->assign('keywords'    , $keywords);    //seo信息
        $this->view->assign('description' , $description); //seo信息

//        $template=$info['template'] ? $info['template'] : ($cate['template_show']?str_replace('.html', '', $cate['template_show']):$this->tablename.'_show');
        $template="";
        if($cate['template_show']){
            $template=str_replace('.html', '', $cate['template_show']);
        }elseif($this->moduledata["template_show"]){
            $template=str_replace('.html', '',$this->moduledata["template_show"]);
        }
        if(empty($template)){
            $template=$this->table.'_show';
        }
//        echo  "详情";
        $template="./template/site_$site_id/$html/$template.html";
        return $this->view->fetch($template);

    }
    public function addonly()
    {
        if (Request::isPost()) {
            $data = Request::post();
            if (empty($data['mobile'])) {
                $this->error("请填写您的联系方式！");
            }
            $data['title'] = "仅提交联系电话";
            $data['createtime'] = time();
            $data['site_id'] = $this->system["id"];
            $data['accountID'] = $this->system["accountID"];
            $data['status'] = 0;
            unset($data['catId']);
            $id = Db::name("website_site_leavemessage")->insertGetId($data);
            if ($id) {
                $this->success("提交成功！");
            } else {
                $this->error("提交失败！");
            }
        } else {
            $this->error("请用POST方式提交数据！");
        }
    }
    //留言表单提交
    public function add(){
        if(Request::isPost()){
            $data = Request::post();
            if(empty($data['name'])){
                $this->error("请填写您的姓名或称呼！");
            }elseif(empty($data['mobile'])){
                $this->error("请填写您的联系方式！");
            }elseif(empty($data['content'])){
                $this->error("请填写留言内容！");
            }
            //是否开启验证码
            $message_code = System::where('id',$this->system["id"])->value('message_code');
            if($message_code){
                if(empty($data['message_code'])){
                    $this->error( '请输入验证码！');
                }elseif( !captcha_check($data['message_code'] )){
                    $this->error( '验证码错误！');
                }else{
                    unset($data['message_code']);
                }
            }
            $data['createtime'] = time();
            $data['site_id'] = $this->system["id"];
            $data['accountID'] = $this->system["accountID"];
            $data['status'] = 0;
            unset($data['catId']);
            $id = Db::name("website_site_leavemessage")->insertGetId($data);
            if($id){
                $this->success("留言成功！");
            }else{
                $this->error("留言失败！");
            }
        }else{
            $this->error("请用POST方式提交数据！");
        }
            /**
            //是否开启验证码
            $message_code = System::where('id',1)
                ->value('message_code');
            if($message_code){
                if( !captcha_check($data['message_code'] ))
                {
                    $result['code']  = '001';
                    $result['msg']  .= '验证码错误;';
                    $this->error($result['msg']);
                }else{
                    unset($data['message_code']);
                }
            }
            //必填项判断
            //查询该模型所有必填字段
            $fields = Db::name('field')
                ->where('moduleid',$this->moduleid)
                ->where('required',1)
                ->field('field,name,errormsg')
                ->select();
            foreach($fields as $k=>$v){
                if(isset($data[$v['field']]) && empty($data[$v['field']]) ){
                    $result['code']  = '001';
                    $result['msg']  = $v['name'].'为必填项';
                }
            }

            if($result['code']  != '001'){
                $id = Db::name("website_site_leavemessage")->insertGetId($data);
                if($id){
                    $result['code']  = '000';
                    $result['msg']  = '留言成功';
                    //邮件通知开始
                    if(System::where('id',1)->value('message_send_mail')){
                        //去除无用字段
                        unset($data['catid']);
                        unset($data['status']);
                        //默认收件人为系统设置中的邮件
                        $email = System::where('id',1)->value('email');
                        $title = '提醒：您的网站有新的留言';
                        //拼接内容
                        $fields = Db::name('field')
                            ->where('moduleid',$this->moduleid)
                            ->field('field,name,type')
                            ->select();
                        $content = '';
                        foreach($fields as $k=>$v){
                            if(isset($data[$v['field']]) ){
                                if($v['type']=='datetime'){
                                    $data[$v['field']] = date("Y-m-d H:i",$data[$v['field']]);
                                }
                                $content .= '<br>'.$v['name'].' : '.$data[$v['field']];
                            }
                        }
                        $this->trySend($email,$title,$content);
                    }
                    //邮件通知结束
                    $this->success($result['msg']);
                }else{
                    $result['code']  = '001';
                    $result['msg']  .= '留言失败;';
                    $this->error($result['msg']);
                }
            }else{
                $this->error($result['msg']);
            }
    */
    }

    //验证码
    public function captcha(){
        $config =    [
            // 验证码字体大小
            'fontSize'    =>    30,
            // 验证码位数
            'length'      =>    4,
            // 关闭验证码杂点
            'useNoise'    =>    true,
            // 是否画混淆曲线
            'useCurve' => false,
        ];
        $captcha = new Captcha($config);
        return $captcha->entry();
    }

    //邮件发送
    private function trySend($email,$title,$content){
        //检查是否邮箱格式
        if (!is_email($email)) {
            return (['code' => 0, 'msg' => '邮箱码格式有误']);
        }
        $send = send_email($email, $title,$content);
        if ($send) {
            return (['code' => 1, 'msg' => '邮件发送成功！']);
        } else {
            return (['code' => 0, 'msg' => '邮件发送失败！']);
        }
    }
    //百度-普通收录
    public function baidu($url=''){
        if(!Cache::get($url)){
            $urls=[];
            if($url){
                array_push($urls,$url);
            }
            $api = $this->system["bd_linksubmit"];
            $ch = curl_init();
            $options =  array(
                CURLOPT_URL => $api,
                CURLOPT_POST => true,
                CURLOPT_RETURNTRANSFER => true,
                CURLOPT_POSTFIELDS => implode("\n", $urls),
                CURLOPT_HTTPHEADER => array('Content-Type: text/plain'),
            );
            curl_setopt_array($ch, $options);
            $result = curl_exec($ch);
            $result_arr=json_decode($result,true);
            if($result_arr["success"]>0){
                Cache::set($url,$result_arr["success"],3600);
            }
            $this->success("收录成功","",$result_arr );
        }else{
            $this->success("刚刚推送收录，过一会再推送");
        }
    }
}
