<?php
/**
 * +----------------------------------------------------------------------
 * | 公共模型基类
 * +----------------------------------------------------------------------
 *                      .::::.
 *                    .::::::::.            | AUTHOR: siyu
 *                    :::::::::::           | EMAIL: 407593529@qq.com
 *                 ..:::::::::::'           | QQ: 407593529
 *             '::::::::::::'               | WECHAT: zhaoyingjie4125
 *                .::::::::::               | DATETIME: 2019/03/04
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
namespace app\common\model;
use think\Model;

class Base extends Model
{
    // 开启自动写入时间戳字段
    protected $autoWriteTimestamp = true;

    // 模型初始化
    protected static function init()
    {
        //TODO:初始化内容
    }

    //通用修改数据
    public function edit($id){
        $info = self::getById($id);
        return $info;
    }

    //通用修改保存
    public function editPost($data)
    {
        $result = self::allowField(true)->save($data, ['id' => $data['id']]);
        if ($result) {
            return ['error' => 0, 'msg' => '修改成功'];
        } else {
            return ['error' => 1, 'msg' => '修改失败'];
        }
    }

    //通用添加保存
    public function addPost($data){
        $result = self::allowField(true)->save($data);
        if ($result) {
            return ['error' => 0, 'msg' => '添加成功'];
        } else {
            return ['error' => 1, 'msg' => '添加失败'];
        }
    }

    //删除
    public function del($id){
        self::destroy($id);
        return ['error'=>0,'msg'=>'删除成功!'];
    }

    //批量删除
    public function selectDel($id){
        self::destroy($id);
        return ['error'=>0,'msg'=>'删除成功!'];
    }

    //排序修改
    public function sort($data){
        $result = self::get($data['id']);
        $result -> sort = $data['sort'];
        $result->save();
        return ['error'=>0,'msg'=>'修改成功!'];

    }

    //状态修改
    public function state($id){
        $data = self::get($id);
        $status = $data['status']==1?0:1;
        $data -> status = $status;
        $data -> save();
        return ['error'=>0,'msg'=>'修改成功!'];
    }


}