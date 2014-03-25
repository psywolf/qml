import QtQuick 2.0

Item {
    width: 320; height: 200

    ListView {
        width: 120;
        model: container.currentModel.len
        delegate: Text {
            text: "I am color number: " + index
            color: container.currentModel.color(index)
        }
        anchors.top: parent.top
        anchors.bottom: parent.bottom
        anchors.horizontalCenter: parent.horizontalCenter
    }
}
