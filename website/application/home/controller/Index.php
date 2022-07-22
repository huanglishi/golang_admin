<?php
/**
 * +----------------------------------------------------------------------
 * | 首页控制器
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


class Index extends Base
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

        //后台开启手机端的时候自动跳转
        if($this->system['ismobile']=='1'){
            if(isMobile()){
                $this->redirect('home/indexmobile/index');
            }
        }
//        echo "后台开启手机端的时候自动跳转：".isMobile();
        $site_id=$this->request->param('site_id');
        $html=$this->system ['html'];
        $this->view->assign('tplappend',"site_$site_id/$html/");//站点id
        $this->view->assign('cate', null);//
        $this->view->assign('system',  $this->system );//系统信息
        $this->view->assign('site_id',     $site_id);//站点id
        if(date("Y",time())>date("Y",$this->system["copyrigstart"])){
            $copyrightdate= $this->system["copyrigstart"]."-".date("Y",time());
        }else{
            $copyrightdate= date("Y",time());
        }
        $this->view->assign('copyrightdate', $copyrightdate);//版权时间
        $this->view->assign('public', '/template/site_'.$this->request->param('site_id').'/');//公共目录
        $this->view->assign('title',  $this->system ['name']);//seo信息
        $this->view->assign('keywords',  $this->system ['keyword']);//seo信息
        $this->view->assign('description',  $this->system ['description']);//seo信息
        $template="./template/site_$site_id/$html/index.html";
//        echo "模板".  $template;
        return $this->view->fetch($template);
    }

    //搜索
    public function search(){
        $search = Request::param('search');//关键字
        $pagesize = Request::param('pagesize')? Request::param('pagesize'):10;//关键字
        if(empty($search)){
            $this->error('请输入关键词');
        }
        $site_id=$this->request->param('site_id');
        $html=$this->system ['html'];
        $this->view->assign('tplappend',"site_$site_id/$html/");//站点id
        $this->view->assign('search', $search);
        $this->view->assign('cate', null);//
        $this->view->assign('system', $this->system );//系统信息
        $this->view->assign('site_id',     $site_id);//站点id
        if(date("Y",time())>date("Y",$this->system["copyrigstart"])){
            $copyrightdate= $this->system["copyrigstart"]."-".date("Y",time());
        }else{
            $copyrightdate= date("Y",time());
        }
        $this->view->assign('copyrightdate', $copyrightdate);//版权时间
        $this->view->assign('public', '/template/site_'.$this->request->param('site_id').'/');//公共目录
        $this->view->assign('title', $this->system['name']);//seo信息
        $this->view->assign('keywords', $this->system['keyword']);//seo信息
        $this->view->assign('description', $this->system['description']);//seo信息
        //分页
        $totalCount = Db::name("website_article_content")
            ->where('title','like',$search)
            ->count();
        $totalPage=ceil($totalCount/$pagesize);
        $pagelist=[];
        for ($item=1; $item<=$totalPage; $item++){
            array_push($pagelist,$item);
        }
        $request = Request::instance();
        $page=Request::param('page');
        if(empty($page)){
            $page=1;
        }
        $pagedata=["baseUrl"=>$request->baseUrl(),"pagelist"=>$pagelist,"totalCount"=>$totalCount,"pagesize"=>$pagesize,
            "totalPage"=>$totalPage,"page"=>$page];
        $this->view->assign('pagedata', $pagedata);//分页
        $template='./template/site_'.$this->request->param('site_id').'/'.$this->system['html'].'/search.html';
        return $this->view->fetch($template);
    }
}
