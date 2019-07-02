import sys
from PyQt5.QtWidgets import QApplication, QWidget, QListWidget, QVBoxLayout
from PyQt5.QtCore import QTimer, Qt
from PyQt5.QtGui import QCursor, QIcon
import ctypes
import requests
import datetime

def getStock(stockName):
    response = requests.get('http://hq.sinajs.cn/list='+stockName)
    listdata = []
    for item in response.text.split(";"):
        if item.isspace():
            continue
        index0 = item.index("\"")
        item = item[index0+1:len(item)-1]
        itemList = item.split(",")
        stock = {}
        stock["name"] = itemList[0]
        oldPrice = float(itemList[2])
        price = float(itemList[3])
        zdfv = float('%.3f' % (price - oldPrice))
        zdf = float('%.2f' % (zdfv / oldPrice * 100))
        stock["oldPrice"] = oldPrice
        stock["price"] = price
        stock["zdfv"] = zdfv
        stock["zdf"] = zdf
        listdata.append(stock)
    return listdata

def getStockList():
    fo = open("stock.list", "r")
    return ",".join(fo.readlines()).replace("\n", "")

class MyApp(QWidget):
    def __init__(self):
        super().__init__()
        self.stockList = getStockList()
        self.resize(200, 30)
        self.setWindowFlags(Qt.FramelessWindowHint|Qt.WindowStaysOnTopHint)
        # self.setWindowOpacity(0.2)
        self.setWindowIcon(QIcon('doraemon.ico'))
        self.setAttribute(Qt.WA_TranslucentBackground)
        self.move(1650, 100)
        self.initUI()
    
    def initUI(self):
        self.list = QListWidget()
        self.list.setStyleSheet("background-color:transparent")
        self.list.setFrameShape(QListWidget.NoFrame)
        self.layout = QVBoxLayout()
        self.layout.addWidget(self.list)
        self.setLayout(self.layout)
        self.first = True
       
        self.stockData()
        self.index = 1
        self.timer = QTimer(self)
        self.timer.timeout.connect(self.stockData)
        self.timer.start(3000)
        self.show()
        
    def stockData(self):
        if not self.checkInTime() and not self.first:
            return
        self.list.clear()
        for item in getStock(self.stockList):
            name = item["name"]
            price = item["price"]
            zdfv = item["zdfv"]
            zdf = item["zdf"]
            str_n = '{0} {1} {2}({3}%)'.format(name, price, zdfv, zdf)
            self.list.addItem(str_n)
        self.first = False
        QApplication.processEvents()
    
    def checkInTime(self):
        amTimeStart = datetime.datetime.strptime(str(datetime.datetime.now().date())+'9:15', '%Y-%m-%d%H:%M')
        amTimeEnd =  datetime.datetime.strptime(str(datetime.datetime.now().date())+'11:30', '%Y-%m-%d%H:%M')
        pmTimeStart = datetime.datetime.strptime(str(datetime.datetime.now().date())+'13:00', '%Y-%m-%d%H:%M')
        pmTimeEnd =  datetime.datetime.strptime(str(datetime.datetime.now().date())+'15:00', '%Y-%m-%d%H:%M')
        nowTime = datetime.datetime.now()
        if amTimeStart <= nowTime and nowTime <= amTimeEnd or pmTimeStart <= nowTime and nowTime <= pmTimeEnd:
            return True
        return False

    def mousePressEvent(self, event):
        if event.button()==Qt.LeftButton:
            self.m_drag=True
            self.m_DragPosition=event.globalPos()-self.pos()
            event.accept()
            self.setCursor(QCursor(Qt.OpenHandCursor))
    def mouseMoveEvent(self, QMouseEvent):
        if Qt.LeftButton and self.m_drag:
            self.move(QMouseEvent.globalPos()-self.m_DragPosition)
            QMouseEvent.accept()
    def mouseReleaseEvent(self, QMouseEvent):
        self.m_drag=False
        self.setCursor(QCursor(Qt.ArrowCursor))

if __name__ == '__main__':
    app = QApplication(sys.argv)
    my = MyApp()
    ctypes.windll.shell32.SetCurrentProcessExplicitAppUserModelID("myappid")
    sys.exit(app.exec_())