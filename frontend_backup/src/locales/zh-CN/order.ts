export default {
  order: {
    list: {
      title: '订单列表',
      search: {
        placeholder: '搜索订单号、客户名称...',
        status: '订单状态',
        dateRange: {
          start: '开始日期',
          end: '结束日期'
        }
      },
      table: {
        orderNumber: '订单号',
        customerName: '客户名称',
        productName: '产品名称',
        quantity: '数量',
        totalPrice: '总价',
        status: '状态',
        createdAt: '创建时间',
        actions: '操作',
        view: '查看',
        edit: '编辑'
      },
      status: {
        pending: '待处理',
        processing: '处理中',
        completed: '已完成',
        cancelled: '已取消'
      }
    },
    detail: {
      title: '订单详情',
      info: {
        title: '订单信息',
        orderNumber: '订单号',
        customerName: '客户名称',
        productName: '产品名称',
        quantity: '数量',
        totalPrice: '总价',
        status: '状态',
        createdAt: '创建时间',
        updatedAt: '更新时间'
      },
      status: {
        title: '状态更新',
        newStatus: '新状态',
        notes: '备注',
        update: '更新状态',
        confirm: {
          title: '确认更新',
          message: '确定要更新订单状态吗？'
        }
      },
      files: {
        title: '文件管理',
        upload: '上传文件',
        download: '下载',
        delete: '删除',
        confirm: {
          title: '确认删除',
          message: '确定要删除此文件吗？'
        },
        tip: '支持上传设计图纸、工艺说明等文件，单个文件不超过10MB'
      },
      actions: {
        edit: '编辑',
        delete: '删除',
        confirm: {
          title: '确认删除',
          message: '确定要删除此订单吗？'
        }
      }
    },
    create: {
      title: '创建订单',
      form: {
        customerName: {
          label: '客户名称',
          placeholder: '请输入客户名称',
          rules: {
            required: '请输入客户名称',
            length: '长度在 2 到 50 个字符'
          }
        },
        productName: {
          label: '产品名称',
          placeholder: '请输入产品名称',
          rules: {
            required: '请输入产品名称',
            length: '长度在 2 到 100 个字符'
          }
        },
        quantity: {
          label: '数量',
          rules: {
            required: '请输入数量',
            min: '数量必须大于0'
          }
        },
        unitPrice: {
          label: '单价',
          rules: {
            required: '请输入单价',
            min: '单价不能为负数'
          }
        },
        totalPrice: {
          label: '总价'
        },
        files: {
          label: '设计文件',
          tip: '支持上传设计图纸、工艺说明等文件，单个文件不超过10MB'
        },
        notes: {
          label: '备注',
          placeholder: '请输入订单备注信息'
        }
      },
      actions: {
        submit: '创建订单',
        cancel: '取消',
        confirm: {
          title: '确认取消',
          message: '确定要取消创建订单吗？'
        }
      }
    },
    messages: {
      create: {
        success: '订单创建成功',
        error: '订单创建失败'
      },
      update: {
        success: '订单更新成功',
        error: '订单更新失败'
      },
      delete: {
        success: '订单删除成功',
        error: '订单删除失败'
      },
      status: {
        success: '状态更新成功',
        error: '状态更新失败'
      },
      file: {
        upload: {
          success: '文件上传成功',
          error: '文件上传失败',
          size: '文件大小不能超过 10MB!'
        },
        download: {
          error: '文件下载失败'
        },
        delete: {
          success: '文件删除成功',
          error: '文件删除失败'
        }
      }
    }
  }
} 