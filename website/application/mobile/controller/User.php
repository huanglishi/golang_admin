<?php
namespace app\mobile\controller;
use think\facade\Request;

class User extends Base
{
    public function initialize()
    {
        parent::initialize();
        //当前模块
        $this->module = strtolower(Request::module());

        $system = db('system')->where('id',1)->find();
        $this->assign('cate', null);//
        $this->assign('system', $system);//系统信息
        $this->assign('public', '/template/'.$this->module.'/'.$system['template'].'/');//公共目录
        $this->assign('title', $system['title']);//seo信息
        $this->assign('keywords', $system['key']);//seo信息
        $this->assign('description', $system['des']);//seo信息
    }
    //用户中心首页
    public function index()
    {
        if(!session('user.id')){
            $this->redirect('login');
        }
        $user = db('users')
            ->alias('u')
            ->leftJoin('users_type ut','u.type_id = ut.id')
            ->field('u.*,ut.name as type_name')
            ->where('u.id',session('user.id'))
            ->find();
        $this->assign('user', $user);
        return view();
    }
    //用户中心设置页
    public function set(){
        if(!session('user.id')){
            $this->redirect('login');
        }
        if(request()->isPost()){
            $data=[];

            //修改密码
            if(input("post.password") && input("post.password2")){
                //密码长度不能低于6位
                if(strlen(trim(input("post.password")))<6){
                    $this->error('密码长度不能低于6位');
                }
                //查看原密码是否正确
                if(input("post.nowpassword")){
                    $id = db('users')
                        ->where('id',session('user.id'))
                        ->where('password',md5(trim(input("post.nowpassword"))))
                        ->find();
                    if(!$id){
                        $this->error('原密码输入有误');
                    }
                }else{
                    $this->error('请输入原密码');
                }
                if(input("post.password") == input("post.password2")){
                    $data['password'] = md5(trim(input("post.password")));
                }else{
                    $this->error('两次输入的密码不一致');
                }
                //更新信息
                db('users')
                    ->where('id', session('user.id'))
                    ->data($data)
                    ->update();
                $this->success('密码修改成功');
            }
            //修改资料
            $data['sex'] = input("post.sex");
            $data['qq'] = input("post.qq");
            $data['mobile'] = input("post.mobile");
            if($data['mobile']){
                //不可和其他用户的一致
                $id = db('users')
                    ->where('mobile',$data['mobile'])
                    ->where('id','<>',session('user.id'))
                    ->find();
                if($id){
                    $this->error('手机号已存在');
                }
            }

            //更新信息
            db('users')
                ->where('id', session('user.id'))
                ->data($data)
                ->update();
            $this->success('修改成功');

        }else{
            $user = db('users')
                ->alias('u')
                ->leftJoin('users_type ut','u.type_id = ut.id')
                ->field('u.*,ut.name as type_name')
                ->where('u.id',session('user.id'))
                ->find();
            $this->assign('user', $user);//
            return view();
        }

    }

    //登录
    public function login(){
        if(request()->isPost()){
            $result=['code'=>'','msg'=>''];
            //登录提交
            $username = trim(input("post.username"));
            $password = trim(input("post.password"));
            //检查是否开启了验证码
            $message_code = db('system')->where('id',1)->value('message_code');
            if($message_code){
                if( !captcha_check(input("post.message_code")))
                {
                    $result['code']  = '001';
                    $result['msg']  .= '验证码错误';
                    $this->error($result['msg']);
                }
            }
            //校验用户名密码
            $user = db('users')
                ->where('email|mobile',$username)
                ->where('password',md5($password))
                ->find();
            if(empty($user)){
                $result['code']  = '001';
                $result['msg']  .= '帐号或密码错误';
                $this->error($result['msg']);
            }else{
                if ($user['status']==1){
                    session('user', $user);
                    //更新信息
                    db('users')
                        ->where('id', $user['id'])
                        ->data(['last_login_time' => time(),'last_login_ip' =>request()->ip()])
                        ->update();

                    $result['code']  = '000';
                    $result['msg']  .= '登录成功';
                    $this->success($result['msg'],'index');
                }else{
                    $result['code']  = '001';
                    $result['msg']  .= '用户已被禁用';
                    $this->error($result['msg']);
                }
            }

        }else{
            if(session('user.id')){
                $this->redirect('index');
            }
            return view();
        }
    }

    //注册
    public function register(){
        if(request()->isPost()){
            $result=['code'=>'','msg'=>''];
            //登录提交
            $email = trim(input("post.email"));
            $password = trim(input("post.password"));
            $password2 = trim(input("post.password2"));

            //密码长度不能低于6位
            if(strlen($password)<6){
                $this->error('密码长度不能低于6位');
            }

            //非空判断
            if(empty($email) || empty($password) || empty($password2)){
                $result['code']  = '001';
                $result['msg']  .= '请输入邮箱、密码和确认密码';
                $this->error($result['msg']);
            }
            //邮箱合法性判断
            if(!is_email($email)){
                $result['code']  = '001';
                $result['msg']  .= '邮箱格式错误';
                $this->error($result['msg']);
            }
            //确认密码
            if($password != $password2){
                $result['code']  = '001';
                $result['msg']  .= '两次密码输入不一致';
                $this->error($result['msg']);
            }

            //检查是否开启了验证码
            $message_code = db('system')->where('id',1)->value('message_code');
            if($message_code){
                if( !captcha_check(input("post.message_code")))
                {
                    $result['code']  = '001';
                    $result['msg']  .= '验证码错误';
                    $this->error($result['msg']);
                }
            }
            //防止重复
            $id = db('users')->where('email|mobile','=',$email)->find();
            if($id){
                $result['code']  = '001';
                $result['msg']  .= '邮箱已被注册';
                $this->error($result['msg']);
            }
            //注册入库
            $data = [];
            $data['email'] = $email;
            $data['password'] = md5($password);
            $data['last_login_time'] = $data['reg_time'] = time();
            $data['reg_ip'] = $data['last_login_ip']=Request::ip();
            $data['status'] = 1;
            $data['type_id'] = 1;
            $id = db('users')->insertGetId($data);
            if($id){
                $this->success('注册成功!','login');
            }else{
                $this->error('注册失败!');
            }
        }else{
            if(session('user.id')){
                $this->redirect('index');
            }
            return view();
        }
    }

    //退出
    public function logout(){
        session('user',null);
        $this->redirect('login');
    }

}
