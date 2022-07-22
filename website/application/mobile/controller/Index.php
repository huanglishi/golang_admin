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
namespace app\mobile\controller;
use app\common\model\System;
use think\Db;
use think\facade\Request;
use think\captcha\Captcha;


class Index extends Base
{
    public function initialize()
    {
        parent::initialize();
        //当前模块
        $this->module = strtolower(Request::module());
    }
    //首页
    public function index()
    {
        $system = System::where('id',1)->find();
        $this->view->assign('cate', null);//
        $this->view->assign('system', $system);//系统信息
        $this->view->assign('public', '/template/'.$this->module.'/'.$system['template'].'/');//公共目录
        $this->view->assign('title', $system['title']);//seo信息
        $this->view->assign('keywords', $system['key']);//seo信息
        $this->view->assign('description', $system['des']);//seo信息
        $template='./template/'.$this->module.'/'.$system['template'].'/'.$system['html'].'/index.html';
        return $this->view->fetch($template);

    }

    //搜索
    public function search(){
        $search = Request::param('search');//关键字
        if(empty($search)){
            $this->error('请输入关键词');
        }
        $this->view->assign('search', $search);

        $system = System::where('id',1)->find();
        $this->view->assign('cate', null);//
        $this->view->assign('system', $system);//系统信息
        $this->view->assign('public', '/template/'.$this->module.'/'.$system['template'].'/');//公共目录
        $this->view->assign('title', $system['title']);//seo信息
        $this->view->assign('keywords', $system['key']);//seo信息
        $this->view->assign('description', $system['des']);//seo信息
        $template='./template/'.$this->module.'/'.$system['template'].'/'.$system['html'].'/search.html';
        return $this->view->fetch($template);
    }

}
