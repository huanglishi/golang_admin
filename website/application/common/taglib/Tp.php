<?php
/**
 * +----------------------------------------------------------------------
 * | 自定义标签
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
namespace app\common\taglib;
use think\template\TagLib;
use think\facade\Request;
class Tp extends TagLib {

    protected $tags = array(
        // 标签定义： attr 属性列表 close 是否闭合（0 或者1 默认1） alias 标签别名 level 嵌套层次
        'close'     => ['attr' => 'time,format', 'close' => 0],                           //闭合标签，默认为不闭合
        'open'      => ['attr' => 'name,type', 'close' => 1],
        'nav'       => ['attr' => 'id,limit', 'close' => 1],                              //通用导航信息
        'cate'      => ['attr' => 'id,type','close' => 0],                                //通用栏目信息
        'catelist'  => ['attr' => 'pid,name,limit','close' => 1],                          //通用栏目列表
        'hascate'   => ['attr' => 'id,name,field','close' => 1],                                //判断是否存在分类
        'position'  => ['attr' => 'name','close' => 1],                                   //通用位置信息
        'link'      => ['attr' => 'name','close' => 1],                                   //获取友情链接
        'ad'        => ['attr' => 'name,type','close' => 1],                              //获取广告信息
        'debris'    => ['attr' => 'name,type','close' => 0],                              //获取碎片信息
        'count'    => ['attr' => 'id','close' => 0],                              // 分类下文章数
        'list'      => ['attr' => 'id,name,pagesize,where,limit,order','close' => 1],     //通用列表
        'search'    => ['attr' => 'search,site_id,table,name,pagesize,where,order','close' => 1], //通用搜索
        'prev'	    => ['attr'	=> 'len','close' => 0],                                   //上一篇
        'prevplus'	=> ['attr'	=> 'name','close' => 1],                                   //上一篇闭环
        'next'	    => ['attr'	=> 'len','close' => 0],                                   //下一篇
        'nextplus'	=> ['attr'	=> 'name','close' => 1],                                   //下一篇闭环

    );

    //这是一个闭合标签的简单演示
    public function tagClose($tag)
    {
        $format = empty($tag['format']) ? 'Y-m-d H:i:s' : $tag['format'];
        $time = empty($tag['time']) ? time() : $tag['time'];
        $parse = '<?php ';
        $parse .= 'echo date("' . $format . '",' . $time . ');';
        $parse .= ' ?>';
        return $parse;
    }
    //获取分类下文章个数
    public function tagCount($tag)
    {
        $id   = $tag['id']?$tag['id']:'';
        $str = '<?php ';
        $str .= 'echo Db::name("website_article_content")->where(\'status\',0)->where("cid",\''.$id.'\')->count();';
        $str .= '?>';
        return $str;
    }

    //这是一个非闭合标签的简单演示
    public function tagOpen($tag, $content)
    {
        $type = empty($tag['type']) ? 0 : 1; // 这个type目的是为了区分类型，一般来源是数据库
        $name = $tag['name']; // name是必填项，这里不做判断了
        $parse = '<?php ';
        $parse .= '$test_arr=[[1,3,5,7,9],[2,4,6,8,10]];'; // 这里是模拟数据
        $parse .= '$__LIST__ = $test_arr[' . $type . '];';
        $parse .= ' ?>';
        $parse .= '{volist name="__LIST__" id="' . $name . '"}';
        $parse .= $content;
        $parse .= '{/volist}';
        return $parse;
    }

    //通用导航信息
    Public function tagNav($tag,$content){
        $tag['limit'] = isset($tag['limit']) ? $tag['limit'] : '0';
        $tag['id']    = isset($tag['id'])    ? $tag['id']    : '';
        $name         = isset($tag['name'])  ? $tag['name']  : 'nav';
        $site_id      =isset($tag['site_id']) ? $tag['site_id'] : '0';

        if(!empty($tag['id'])){
            $catestr = '$__CATE__ = Db::name(\'website_article_cate\')->where(\'is_menu\',1)->order(\'weigh ASC,id DESC\')->select();';
            $catestr.= '$__LIST__ = getChildsOn($__CATE__,'.$tag['id'].');';
        }else{
            $catestr = '$__CATE__ = Db::name(\'website_article_cate\')->where(\'is_menu\',1)->where(\'site_id\','.$site_id.')->order(\'weigh ASC,id DESC\')->select();';
            $catestr.= '$__LIST__ = unlimitedForLayer($__CATE__);';
        }
        //提取前N条数据,因为sql的LIMIT避免不了子栏目的问题
        if(!empty($tag['limit'])){
            $catestr.= '$__LIST__ =  array_slice($__LIST__, 0,'.$tag['limit'].');';
        }
        $parse = '<?php ';
        $parse .= $catestr;
        $parse .= ' ?>';
        $parse .= '{volist name="__LIST__" id="' . $name . '"}';
        $parse .= $content;
        $parse .= '{/volist}';
        return $parse;
    }

    //通用栏目信息
    Public function tagCate($tag){
        $id   = isset($tag['id'])?$tag['id']:"input('catId')";
        $type = $tag['type']?$tag['type']:'catname';

        $str = '<?php ';
        $str .= '$__CATE__=Db::name("website_article_cate")->where("id",'.$id.')->find();';
        $str .= 'if(is_array($__CATE__)){ ';
        $str .= '$__CATE__[\'url\']=getUrl($__CATE__);';
        $str .= 'echo $__CATE__[\''.$type.'\'];';
        $str .= '}';
        $str .= '?>';
        return $str;
    }

    //通用位置信息
    Public function tagPosition($tag,$content){
        $name   = $tag['name']?$tag['name']:'position';
        $parse  = '<?php ';
        $parse .= '$__CATE__ = Db::name(\'website_article_cate\')->select();';
        $parse .= '$__LIST__ = getParents($__CATE__,input(\'catId\'));';
        $parse .= ' ?>';
        $parse .= '{volist name="__LIST__" id="' . $name . '"}';
        $parse .= '<?php $' . $name . '[\'url\']=getUrl( $' . $name . '); ?>';
        $parse .= $content;
        $parse .= '{/volist}';
        return $parse;
    }

    //获取友情链接
    Public function tagLink($tag,$content){
        $name   = $tag['name']?$tag['name']:'link';
        $site_id      =isset($tag['site_id']) ? $tag['site_id'] : '0';
        $parse  = '<?php ';
        $parse .= '$__LIST__ = Db::name(\'website_link\')->where(\'status\',1)->where(\'site_id\','.$site_id.')->order(\'weigh ASC,id desc\')->select();';
        $parse .= ' ?>';
        $parse .= '{volist name="__LIST__" id="' . $name . '"}';
        $parse .= $content;
        $parse .= '{/volist}';
        return $parse;
    }

    //获取广告信息
    Public function tagAd($tag,$content){
        $name = isset($tag['name']) ? $tag['name'] : 'ad';
        $flag = isset($tag['flag']) ? $tag['flag'] : '';
        $cid   = isset($tag['cid'])   ? $tag['cid']   : '';
        $parse = '<?php ';
        $parse .= '
            $__WHERE__ = array();
            if (!empty(\'' . $cid . '\')) {
                $__WHERE__[] = [\'cid\', \'=\', ' . $cid . '];
            }
            if (!empty(\'' . $flag . '\')) {
                $__WHERE__[] = [\'flag\', \'like\', \'%' . $flag . '%\'];
            }';
        $parse .= '$__LIST__ = Db::name(\'website_article_content\')
            ->field(\'id,title,swipimg,image,des as description,link as url\')
            ->where(\'status\',0)
            ->where(\'site_id\',input(\'site_id\'))
            ->where($__WHERE__)
            ->order(\'weigh ASC,id desc\')
            ->select();
            $__LIST__=changImg($__LIST__);
            ';
        $parse .= ' ?>';
        $parse .= '{volist name="__LIST__" id="' . $name . '"}';
        $parse .= $content;
        $parse .= '{/volist}';
        return $parse;
    }

    //通用碎片信息
    Public function tagDebris($tag){
        $id   = $tag['id']?$tag['id']:'';
        $type   = $tag['type']?$tag['type']:'';
        $str = '<?php ';
        $str .= 'echo Db::name("website_article_content")->where("id",\''.$id.'\')->value("'.$type.'");';
        $str .= '?>';
        return $str;
    }

    //通用列表
    Public function tagList($tag,$content){
        $id    = isset($tag['id'])    ? $tag['id']     : "input('catId')";                //可以为空
        $name  = isset($tag['name'])  ?  $tag['name']  : "list";                          //不可为空
        $order = isset($tag['order']) ?  $tag['order'] : 'weigh ASC,id DESC';              //排序
        $limit = isset($tag['limit']) ?  $tag['limit'] : '0';                             //多少条数据,传递时不再进行分页
        $where = isset($tag['where']) ?  $tag['where'].' AND status = 0 ' : 'status = 0'; //查询条件
        $pagesize = isset($tag['pagesize']) ?  $tag['pagesize']   : config('page_size');
        //paginate(1,false,['query' => request()->param()]); //用于传递所有参数，目前只需要page参数

        $parse  = '<?php ';
        $parse .='
            //查找栏目对应的表信息
            $__TABLE_=Db::name(\'website_article_cate\')
                ->alias(\'a\')
                ->leftJoin(\'website_module m\',\'a.module_id = m.id\')
                ->field(\'a.id,a.module_id,a.pagesize,a.catname,a.site_id,m.model_name as modulename\')
                ->where(\'a.id\',\'=\','.$id.')
                ->find();
            //获取表名称    
            $__TABLENAME_ = $__TABLE_[\'modulename\'];
            //获取模型ID
            $__MODULEID__ = $__TABLE_[\'module_id\'];
            $__SITEID__ = $__TABLE_[\'site_id\'];
            //查询子分类,列表要包含子分类内容
            $__ALLCATE__ = Db::name(\'website_article_cate\')->field(\'id,pid\')->select();
            $__IDS__ = getChildsIdStr(getChildsId($__ALLCATE__,'.$id.'),'.$id.');

            //表名称为空时（id不存在）直接返回
            if(!empty($__TABLENAME_)){
                //当传递limit时，不再进行分页
                if('.$limit.'!=0){
                    $__LIST__ = Db::name($__TABLENAME_)
                    ->order(\''.$order.'\')
                    ->limit(\''.$limit.'\')
                    ->where(" '.$where.'")
                    ->where(\'cid\',\'in\',$__IDS__)
                    ->select();
                    $page = \'\';
                }else{
                    $__TABLE_[\'pagesize\'] = empty($__TABLE_[\'pagesize\'])?'.$pagesize.':$__TABLE_[\'pagesize\'];
                    $__LIST__ = Db::name($__TABLENAME_)
                    ->order(\''.$order.'\')
                    ->where(" '.$where.'")
                    ->where(\'cid\',\'in\',$__IDS__)
                    ->paginate($__TABLE_[\'pagesize\']);
                    $page = $__LIST__->render();
                }
                //处理数据（把列表中需要处理的字段转换成数组和对应的值）
                $__LIST__ = changeFields($__LIST__,$__SITEID__);
            }else{
                $__LIST__ = [];
            }
            ';
        $parse .= ' ?>';
        $parse .= '{volist name="__LIST__" id="'.$name.'"}';
        $parse .= $content;
        $parse .= '{/volist}';
        return $parse;
    }

    //通用搜索 search,table,name,pagesize,where,order
    Public function tagSearch($tag,$content){
        $search = isset($tag['search'])    ? $tag['search']     : "";                     //关键字
        $site_id = isset($tag['site_id'])    ? $tag['site_id']     : "";                     //关键字
        $table  = isset($tag['table'])     ? $tag['table']      : "website_article_content";              //表名称
        $name  = isset($tag['name'])  ?  $tag['name']  : "list";                          //不可为空
        $order = isset($tag['order']) ?  $tag['order'] : 'weigh ASC,id DESC';              //排序
        $where = isset($tag['where']) ?  $tag['where'].' AND status = 0 ' : 'status = 0'; //查询条件
        $pagesize = isset($tag['pagesize']) ?  $tag['pagesize']   : config('page_size');

        $parse  = '<?php ';
        $parse .='
                $__LIST__ = Db::name("'.$table.'")
                ->order("'.$order.'")
                ->where(\'site_id\',input(\'site_id\'))
                ->where("'.$where.'")
                ->where("title", "like", "%'.$search.'%")
                ->paginate("'.$pagesize.'",false,[\'query\' => request()->param()]);
                $page = $__LIST__->render();

            //处理数据（把列表中需要处理的字段转换成数组和对应的值）
            $__LIST__ = changeFieldsSearch($__LIST__,'.$site_id.');
            ';
        $parse .= ' ?>';
        $parse .= '{volist name="__LIST__" id="'.$name.'"}';
        $parse .= $content;
        $parse .= '{/volist}';
        return $parse;
    }

    //详情上一篇
    Public function tagPrev($tag){
        $len = $tag['len']?$tag['len']:'';

        $str  = '<?php ';
        $str .= '
                //查找表名称
                $__TABLENAME__ = Db::name(\'website_article_cate\')
                    ->alias(\'c\')
                    ->leftJoin(\'website_module m\',\'c.module_id = m.id\')
                    ->field(\'m.model_name as table_name\')
                    ->where(\'c.id\',input(\'catId\'))
                    ->find();
                //根据ID查找上一篇的信息
                $__PREV__ = Db::name($__TABLENAME__[\'table_name\'])
                    ->where(\'cid\',input(\'catId\'))
                    ->where(\'id\',\'<\',input(\'id\'))
                    ->field(\'id,cid,title\')
                    ->order(\'weigh ASC,id DESC\')
                    ->find();
                if($__PREV__){
                    //处理上一篇中的URL
                    $__PREV__[\'url\'] = getShowUrl($__PREV__);
                    $__PREV__ = "<a class=\"prev\" title=\" ".$__PREV__[\'title\']." \" href=\" ".$__PREV__[\'url\']." \">".$__PREV__[\'title\']."</a>"; 
                }else{
                    $__PREV__ = "<a class=\"prev_no\" href=\"javascript:;\">暂无数据</a>"; 
                }
                echo $__PREV__;
                ';
        $str .= '?>';
        return $str;
    }

    //详情下一篇
    Public function tagNext($tag){
        $len = $tag['len']?$tag['len']:'';

        $str  = '<?php ';
        $str .= '
                //查找表名称
                $__TABLENAME__ = Db::name(\'website_article_cate\')
                    ->alias(\'c\')
                    ->leftJoin(\'website_module m\',\'c.module_id = m.id\')
                    ->field(\'m.model_name as table_name\')
                    ->where(\'c.id\',input(\'catId\'))
                    ->find();
                //根据ID查找下一篇的信息
                $__PREV__ = Db::name($__TABLENAME__[\'table_name\'])
                    ->where(\'cid\',input(\'catId\'))
                    ->where(\'id\',\'>\',input(\'id\'))
                    ->field(\'id,cid,title\')
                    ->order(\'weigh ASC,id DESC\')
                    ->find();
                if($__PREV__){
                    //处理下一篇中的URL
                    $__PREV__[\'url\'] = getShowUrl($__PREV__);
                    $__PREV__ = "<a class=\"next\" title=\" ".$__PREV__[\'title\']." \" href=\" ".$__PREV__[\'url\']." \">".$__PREV__[\'title\']."</a>"; 
                }else{
                    $__PREV__ = "<a class=\"next_no\" href=\"javascript:;\">暂无数据</a>"; 
                }
                echo $__PREV__;
                ';
        $str .= '?>';
        return $str;
    }
/*****翻篇plus**/
    //详情上一篇
    public function tagPrevplus($tag,$content){
        $name   = $tag['name']?$tag['name']:'item';
        $parse  = '<?php ';
        $parse .= '
            //查找表名称
             $__TABLENAME__ = Db::name(\'website_article_cate\')
                    ->alias(\'c\')
                    ->leftJoin(\'website_module m\',\'c.module_id = m.id\')
                    ->field(\'m.model_name as table_name\')
                    ->where(\'c.id\',input(\'catId\'))
                    ->find();
              //根据ID查找上一篇的信息
            $__LIST__ =Db::name($__TABLENAME__[\'table_name\'])
                    ->where(\'cid\',input(\'catId\'))
                    ->where(\'id\',\'<\',input(\'id\'))
                    ->limit(1)
                    ->field(\'id,cid,title,image,createtime\')
                    ->order(\'weigh ASC,id DESC\')->select();';
        $parse .= ' ?>';
        $parse .= '{volist name="__LIST__" id="' . $name . '"}';
        $parse .= '<?php $' . $name . '[\'url\']=getShowUrl( $' . $name . '); ?>';
        $parse .= $content;
        $parse .= '{/volist}';
        return $parse;
    }
    //详情下一篇-闭环
    public function tagNextplus($tag,$content){
        $name   = $tag['name']?$tag['name']:'item';
        $parse  = '<?php ';
        $parse .= '
          //查找表名称
              //查找表名称
                $__TABLENAME__ = Db::name(\'website_article_cate\')
                    ->alias(\'c\')
                    ->leftJoin(\'website_module m\',\'c.module_id = m.id\')
                    ->field(\'m.model_name as table_name\')
                    ->where(\'c.id\',input(\'catId\'))
                    ->find();
                //根据ID查找下一篇的信息
                $__LIST__ = Db::name($__TABLENAME__[\'table_name\'])
                    ->where(\'cid\',input(\'catId\'))
                    ->where(\'id\',\'>\',input(\'id\'))
                    ->field(\'id,cid,title,image,createtime\')
                   ->limit(1)
                    ->order(\'weigh ASC,id DESC\')->select();';
        $parse .= ' ?>';
        $parse .= '{volist name="__LIST__" id="' . $name . '"}';
        $parse .= '<?php $' . $name . '[\'url\']=getShowUrl( $' . $name . '); ?>';
        $parse .= $content;
        $parse .= '{/volist}';
        return $parse;
    }
    //获取分类列表
    Public function tagCatelist($tag,$content){
        $pid    = isset($tag['pid'])    ? $tag['pid'] : '';
        $name   = isset($tag['name'])  ? $tag['name']  : 'list';
        $limit = isset($tag['limit']) ?  $tag['limit'] : '10';                             //多少条数据
        $parse  = '<?php ';
        $parse .= '
                $__LIST__ = Db::name(\'website_article_cate\')
               ->where("pid",'.$pid.')
               ->limit(\''.$limit.'\')
                ->field(\'id,catname,url,pagesize,is_next,is_blank,catdir,module_id,icoimage\')
                ->order(\'weigh ASC,id DESC\')->select();';
        $parse .= ' ?>';
        $parse .= '{volist name="__LIST__" id="' . $name . '"}';
        $parse .= '<?php $' . $name . '[\'url\']=getUrl( $' . $name . '); ?>';
        $parse .= $content;
        $parse .= '{/volist}';
        return $parse;
    }
    //是否存在分类
    Public function tagHascate($tag,$content){
        $id    = isset($tag['id'])    ? $tag['id'] : '';
        $field    = isset($tag['field'])    ? $tag['field'] : '';
        if($field){
            $field="id,catname,url,pagesize,is_next,is_blank,catdir,module_id,".$field;
        }else{
            $field="id,catname,url,pagesize,is_next,is_blank,catdir,module_id";
        }
        $name   = isset($tag['name'])  ? $tag['name']  : 'list';
        $parse  = '<?php ';
        $parse .= '
                $__LIST__ = Db::name(\'website_article_cate\')
                ->where("id",'.$id.')
                ->where("status",0)
                ->field("'.$field.'")
                ->order(\'weigh ASC,id DESC\')->select();';
        $parse .= ' ?>';
        $parse .= '{volist name="__LIST__" id="' . $name . '"}';
        $parse .= '<?php $' . $name . '[\'url\']=getUrl( $' . $name . '); ?>';
        $parse .= $content;
        $parse .= '{/volist}';
        return $parse;
    }

}