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
namespace app\mobile\controller;
use app\common\model\Cate;
use app\common\model\Module as M;
use app\common\model\System;
use think\Db;
use think\facade\Request;
use think\captcha\Captcha;

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
        $this->moduleid = Cate::where('id','=',input('catId'))
            ->value('moduleid');
        //当前表
        $this->tablename = M::where('id','=',$this->moduleid)
            ->value('name');
    }
    //列表
    public function index(){
        if(Request::param('catId')){
            $cate = Cate::where('id','=',Request::param('catId'))
                ->find();
            $cate['topid'] = $cate['parentid']?$cate['parentid']:$cate['id'];
        }else{
            return redirect("/");
        }
        $system = System::where('id',1)
            ->find();

        if($this->tablename=='page'){
            //单页模型
            //查找第一条记录
            $info = Db::name($this->tablename)
                ->where('catid','=',Request::param('catId'))
                ->find();
            $this->view->assign('info', $info);//单页内容
            //定义tdk
            $title  = $cate['title'] ? $cate['title'] :
                ($cate['catname'] ? $cate['catname'] : $info['title']);
            $keywords    = $cate['keywords'] ? $cate['keywords'] : $system['key'];
            $description = $cate['description'] ? $cate['description'] : $system['des'];
        }else{
            //列表模型
            //定义tdk
            $title       = $cate['title']       ? $cate['title']       : $cate['catname'];
            $keywords    = $cate['keywords']    ? $cate['keywords']    : $system['key'];
            $description = $cate['description'] ? $cate['description'] : $system['des'];
        }
        $this->view->assign('cate',        $cate);//栏目信息
        $this->view->assign('system',      $system);//系统信息
        $this->view->assign('public',      '/template/'.$this->module.'/'.$system['template'].'/');//公共目录
        $this->view->assign('title',       $title);//seo信息
        $this->view->assign('keywords',    $keywords);//seo信息
        $this->view->assign('description', $description);//seo信息

        $template=$cate['template_list']?str_replace('.html', '', $cate['template_list']):$this->tablename.'_list';
        return $this->view->fetch($template);

    }

    //详情
    public function info(){
        if(Request::param('id')){
            //点击数增加
            Db::name($this->tablename)
                ->where('id',Request::param('id'))
                ->setInc('hits');
        }
        if(Request::param('catId')){
            $cate = Cate::where('id','=',Request::param('catId'))->find();
            $cate['topid'] = $cate['parentid']?$cate['parentid']:$cate['id'];
        }
        $system = System::where('id',1)
            ->find();
        //查找详情信息
        $info = Db::name($this->tablename)
            ->where('id',input('id'))
            ->find();
        $info = changefield($info,$this->moduleid);
        $info['url'] = getShowUrl($info);

        //定义tdk
        $title  = $cate['title'] ? $cate['title'] : $cate['catname'];
        $keywords = $info['keywords']?$info['keywords']:($cate['keywords'] ? $cate['keywords'] : $system['key']);
        $description = $info['description']?$info['description']:($cate['description'] ? $cate['description'] : $system['des']);

        $this->view->assign('cate'        , $cate);        //栏目信息
        $this->view->assign('info'        , $info);        //详情信息
        $this->view->assign('system'      , $system);      //系统信息
        $this->view->assign('public'      , '/template/'.$this->module.'/'.$system['template'].'/');//公共目录
        $this->view->assign('title'       , $title);       //seo信息
        $this->view->assign('keywords'    , $keywords);    //seo信息
        $this->view->assign('description' , $description); //seo信息

        $template=$info['template'] ? $info['template'] : ($cate['template_show']?str_replace('.html', '', $cate['template_show']):$this->tablename.'_show');
        return $this->view->fetch($template);

    }

    //留言表单提交
    public function add(){
        $result=['code'=>'','msg'=>''];
        if(Request::isPost()){
            $data = Request::post();
            $data['create_time'] = time();
            $data['catid'] = $data['catId'];
            $data['status'] = 0;
            unset($data['catId']);

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
                $id = Db::name($this->tablename)->insertGetId($data);
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
                        $title = 'SIYUCMS提醒：您的网站有新的留言';
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

        }
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

}
