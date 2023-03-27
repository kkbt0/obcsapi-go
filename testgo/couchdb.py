import couchdb
import time

couch = couchdb.Server('http://admin:password@localhost:5984/')

db = couch['obs']

# for id in db:
#  print(id)

def append(filename: str,new_str: str):
        """用于纯文本 判断文件是否存在 不存在则新建，存在则结尾添加"""
        timestmp13 = int(time.time() * 1000)
        new_node_id = "h:" + str(timestmp13)
        if filename in db:
                print("File exists")
                file_json = db[filename]
                # print(file_json['children'][-1]) # 取出最后一个元素

                dt_child = db[file_json['children'][-1]]
                new_data = dt_child['data'] + new_str
                if(file_json.get('deleted')==True):
                        new_data = new_str # 清空原来的
                        file_json['deleted'] = False
                
                new_node = {
                "_id": new_node_id,
                "data": new_data,
                "type": "leaf"
                }
                print(db.save(new_node))

                file_json['mtime'] =int(time.time() * 1000)
                file_json['children'] = [new_node_id]
                file_json['size'] = len(new_data)
                print(db.update([file_json]))
        else:
                print("File Not exists")
                new_node = {
                "_id": new_node_id,
                "data": new_str,
                "type": "leaf"
                }
                print(db.save(new_node))

                new_file_node = {
                "_id": filename,
                "children": [new_node_id],
                "ctime": timestmp13,
                "mtime": timestmp13,
                "size": len(new_data),
                "type": "plain"
                }
                print(db.save(new_file_node))

def read_file(filepath):
        with open(filepath, 'rb') as f1:
                base64_str = base64.b64encode(f1.read())
                return base64_str
        
import base64

def store_file(filename: str,new_byte64):
        """用于存储文件、图片"""
        timestmp13 = int(time.time() * 1000)
        new_node_id = "h:" + str(timestmp13)
        new_node = {
                "_id": new_node_id,
                "data": new_byte64,
                "type": "leaf"
                }
        print(db.save(new_node)) # 存文件
        if filename in db:
                print("File exists")
                file_json = db[filename]
                # print(file_json['children'][-1]) # 取出最后一个元素
                 # True False None
                if(file_json.get('deleted')==True):
                        file_json['deleted'] = False
                
                file_json['mtime'] =int(time.time() * 1000)
                file_json['children'] = [new_node_id]
                file_json['size'] = len(new_byte64)
                print(db.update([file_json]))
        else:
                print("File Not exists")
                new_file_node = {
                "_id": filename,
                "children": ["h:8r6xx3"],
                "ctime": timestmp13,
                "mtime": timestmp13,
                "size": len(new_byte64),
                "type": "newnote"
                }
                print(db.save(new_file_node))


append("2.md","new data\n")
store_file("test0.jpg",read_file("test.jpg"))

#print(read_file("test.jpg"))