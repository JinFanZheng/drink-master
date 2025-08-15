# GORM模型迁移计划

## 概述
基于生产数据库表结构分析，对现有GORM模型进行完整的结构调整，确保严格匹配数据库表结构。

## 数据库表结构分析结果

### 核心表结构
1. **members** - 用户信息表
2. **products** - 产品信息表  
3. **orders** - 订单表
4. **machines** - 机器信息表
5. **machine_owners** - 机器所有者表
6. **franchise_intentions** - 加盟意向表
7. **material_silos** - 物料槽表
8. **machine_product_prices** - 机器产品价格表

## 模型变更详细计划

### ✅ 1. Member模型 - 已正确匹配
- **状态**: 完全匹配数据库结构
- **字段对比**: 所有字段类型和指针类型都正确
- **无需更改**

### ✅ 2. Product模型 - 已正确匹配  
- **状态**: 完全匹配数据库结构
- **字段对比**: 包括Price和PriceWithoutCup字段
- **无需更改**

### ✅ 3. Order模型 - 已正确匹配
- **状态**: 完全匹配数据库结构
- **字段对比**: 所有指针类型都正确处理
- **无需更改**

### ✅ 4. MachineOwner模型 - 已正确匹配
- **状态**: 完全匹配数据库结构
- **无需更改**

### ✅ 5. FranchiseIntention模型 - 已正确匹配
- **状态**: 完全匹配数据库结构  
- **无需更改**

### ❌ 6. Machine模型 - 需要重大调整

**数据库字段结构**:
```sql
Id: varchar(36) NO PRI 
MachineOwnerId: varchar(36) YES  
MachineNo: varchar(32) YES  
Name: varchar(32) YES  
Area: varchar(64) YES  
Address: varchar(128) YES  
ServicePhone: varchar(11) YES  
BusinessStatus: int NO  
SubscribeTime: datetime(3) YES  
UnSubscribeTime: datetime(3) YES  
DeviceId: varchar(255) YES  
DeviceName: varchar(255) YES  
DeviceSn: varchar(255) YES  
BindDeviceTime: datetime(3) YES  
Version: bigint NO  
CreatedOn: datetime(3) NO  
UpdatedOn: datetime(3) YES  
IsDebugMode: bit(1) NO  
```

**当前模型问题**:
1. `MachineOwnerId` 应该是 `*string` (可空)
2. `MachineNo` 应该是 `*string` (可空)  
3. `Name` 应该是 `*string` (可空)
4. `Area` 应该是 `*string` (可空)
5. `Address` 应该是 `*string` (可空)
6. `IsDebugMode` 应该使用 `bool` 类型，不是 `[]byte`

### ❌ 7. MaterialSilo模型 - 需要完全重构

**数据库字段结构**:
```sql
Id: varchar(36) NO PRI 
MachineId: varchar(36) YES  
No: varchar(16) YES  
Type: int NO  
ProductId: varchar(255) YES  
IsSale: bit(1) NO  
Total: int NO  
Stock: int NO  
SingleFeed: int NO  
Version: bigint NO  
CreatedOn: datetime(3) NO  
UpdatedOn: datetime(3) YES  
```

**当前模型问题**:
1. 字段名称不匹配：`SiloNo` vs `No`
2. 缺少 `Type` 字段
3. 缺少 `IsSale` 字段 
4. 缺少 `Total` 字段
5. 缺少 `SingleFeed` 字段
6. 使用错误的时间字段名称
7. 字段类型不匹配

### ✅ 8. MachineProductPrice模型 - 已正确匹配
- **状态**: 完全匹配数据库结构
- **无需更改**

## 实施优先级

### 高优先级 (必须修复)
1. **Machine模型字段类型修正**
2. **MaterialSilo模型完全重构**

### 中优先级  
3. **测试文件更新** - 修复所有指针类型相关的测试错误

### 低优先级
4. **代码优化** - 清理unused imports和注释

## 风险评估

### 🔴 高风险变更
- MaterialSilo模型重构：影响物料管理相关功能
- Machine模型字段类型变更：影响机器管理功能

### 🟡 中风险变更  
- 测试文件更新：可能影响CI/CD流程

### 🟢 低风险变更
- 已匹配的模型：无需更改

## 测试策略

### 1. 单元测试
- 每个模型变更后运行对应的单元测试
- 验证字段映射和类型转换

### 2. 集成测试
- 验证数据库连接和CRUD操作
- 验证模型关联和查询

### 3. API测试
- 验证所有相关API接口正常工作
- 验证JSON序列化/反序列化

## 实施时间表

### 阶段1: 模型修复 (优先级：高)
- [ ] 修复Machine模型字段类型
- [ ] 重构MaterialSilo模型
- [ ] 验证模型与数据库结构匹配

### 阶段2: 测试更新 (优先级：中)  
- [ ] 更新所有相关测试文件
- [ ] 修复指针类型转换问题
- [ ] 运行完整测试套件

### 阶段3: 验证 (优先级：高)
- [ ] 数据库连接测试
- [ ] API功能测试  
- [ ] 生产环境兼容性验证

## 备注

- 所有变更必须保持与生产数据库的严格一致性
- 禁用模型关联以避免复杂的字段映射问题
- 使用手动SQL查询替代GORM关联查询
- 保持向后兼容性，避免breaking changes