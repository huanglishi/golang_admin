<?php
/**
 * +----------------------------------------------------------------------
 * | 首页控制器-手机
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
use app\common\model\System;
use think\Db;
use think\facade\Request;
use think\captcha\Captcha;


class Indexmobile extends Base
{
    private $system=null;
    public function initialize()
    {
        parent::initialize();
        $this->system = System::where('id',$this->request->param('site_id'))->find();
        if($this->system['wxqrcode']){
            $this->system['wxqrcode']= $this->system['file_root']. $this->system['wxqrcode'];
        }
        if($this->system['logo']){
            $this->system['logo']= $this->system['file_root']. $this->system['logo'];
        }
    }
    //首页
    public function index()
    {

        //手机版跳 pc
        if(!isMobile()){
            $this->redirect('home/index/index');
        }
//        echo "后台开启手机端的时候自动跳转：".isMobile();
        $this->view->assign('cate', null);//
        $this->view->assign('system',  $this->system );//系统信息
        $this->view->assign('public', '/template/site_'.$this->request->param('site_id').'/');//公共目录
        $this->view->assign('title',  $this->system ['name']);//seo信息
        $this->view->assign('keywords',  $this->system ['keyword']);//seo信息
        $this->view->assign('description',  $this->system ['description']);//seo信息
        $template='./template/site_'.$this->request->param('site_id').'/'. $this->system['mobile_html'].'/index.html';
        return $this->view->fetch($template);
    }

    //搜索
    public function search(){
        $search = Request::param('search');//关键字
        if(empty($search)){
            $this->error('请输入关键词');
        }
        $this->view->assign('search', $search);
        $this->view->assign('cate', null);//
        $this->view->assign('system', $this->system );//系统信息
        $this->view->assign('public', '/template/site_'.$this->request->param('site_id').'/');//公共目录
        $this->view->assign('title', $this->system['name']);//seo信息
        $this->view->assign('keywords', $this->system['keyword']);//seo信息
        $this->view->assign('description', $this->system['description']);//seo信息
        $template='./template/site_'.$this->request->param('site_id').'/'.$this->system['mobile_html'].'/search.html';
        return $this->view->fetch($template);
    }

}
