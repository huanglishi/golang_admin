<?php
/**
 * +----------------------------------------------------------------------
 * | 路由部分
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
//完整域名绑定到模块
//前台路由部分
$sitelist= Db::name('website_site')
    ->field('id,domain')
    ->where("status",0)
    ->where("step",1)
    ->select();
foreach ($sitelist as$site){
   if($site["domain"]){
       Route::domain($site["domain"], 'home')->append(['site_id'=>$site["id"]]);
   }
}
//Route::domain('ynjiyuan.com', 'home')->append(['site_id'=>1]);
//1.1公司官网
//Route::domain('web.linbint.com', 'home')->append(['site_id'=>1]);
////1.1.2公司官网新版测试
//Route::domain('www.linbint.com', 'home')->append(['site_id'=>2]);
////1.2-1云南风云拾光有限司
//Route::domain('ts.linbint.com', 'home')->append(['site_id'=>4]);
////1.2-2云南风云拾光有限司
//Route::domain('cm.linbint.com', 'home')->append(['site_id'=>6]);
////1.3云南硕瑜环保科技有限公司
//Route::domain('vsite.linbint.com', 'home')->append(['site_id'=>5]);
//前台路由部分
$cate = Db::name('website_article_cate')
    ->alias('a')
    ->leftJoin('website_module m','a.module_id = m.id')
    ->field('a.id,a.catname,a.catdir,m.name as modulename,m.model_name as moduleurl')
    ->order('a.weigh ASC,a.id ASC')
    ->select();
$home_rote=[];
foreach ($cate as $k=>$v){
    //只有设置了栏目目录的栏目才配置路由
    if($v['catdir']){
        if($v['moduleurl']=='page'){
            //单页模型
            //1.1公司官网
            $home_rote[''.$v['catdir'].'-:catId'] = 'home/'.$v['catdir'].'/index';
            //Mobile
            $home_rote['mobile/'.$v['catdir'].'-:catId'] = 'mobile/'.$v['catdir'].'/index';
        }else{
            //列表+详情模型
            //1.1公司官网
            $home_rote[''.$v['catdir'].'-:catId/:id'] = 'home/'.$v['catdir'].'/info';
            $home_rote[''.$v['catdir'].'-:catId'] = 'home/'.$v['catdir'].'/index';
            //Mobile
            $home_rote['mobile/'.$v['catdir'].'-:catId/:id'] = 'mobile/'.$v['catdir'].'/info';
            $home_rote['mobile/'.$v['catdir'].'-:catId'] = 'mobile/'.$v['catdir'].'/index';
        }
    }
}
return $home_rote;

