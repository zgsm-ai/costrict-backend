#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
#include <memory>
#include <map>

enum class DeviceType {
    LIGHT,
    THERMOSTAT,
    DOOR_LOCK
};

enum class DeviceStatus {
    ON,
    OFF,
    STANDBY,
    ERROR
};

class SmartDevice {
protected:
    std::string deviceId;
    std::string name;
    DeviceType type;
    DeviceStatus status;
    std::string location;
    
public:
    SmartDevice(const std::string& id, const std::string& name,
                DeviceType type, const std::string& location)
        : deviceId(id), name(name), type(type), status(DeviceStatus::OFF),
          location(location) {}
    
    virtual ~SmartDevice() = default;
    
    std::string getId() const { return deviceId; }
    std::string getName() const { return name; }
    DeviceType getType() const { return type; }
    DeviceStatus getStatus() const { return status; }
    std::string getLocation() const { return location; }
    
    void setName(const std::string& newName) { name = newName; }
    void setLocation(const std::string& newLocation) { location = newLocation; }
    
    virtual bool turnOn() {
        status = DeviceStatus::ON;
        std::cout << name << " ON" << std::endl;
        return true;
    }
    
    virtual bool turnOff() {
        status = DeviceStatus::OFF;
        std::cout << name << " OFF" << std::endl;
        return true;
    }
    
    virtual void getStatusReport() const {
        std::cout << name << " (" << deviceId << ") - " << location << " - " << static_cast<int>(status) << std::endl;
    }
};

class SmartLight : public SmartDevice {
private:
    int brightness;
    std::string color;
    
public:
    SmartLight(const std::string& id, const std::string& name,
               const std::string& location, int brightness = 100,
               const std::string& color = "#FFFFFF")
        : SmartDevice(id, name, DeviceType::LIGHT, location),
          brightness(brightness), color(color) {}
    
    void setBrightness(int level) {
        brightness = std::max(0, std::min(100, level));
    }
    
    void setColor(const std::string& hexColor) {
        color = hexColor;
    }
    
    void getStatusReport() const override {
        std::cout << name << " - Brightness: " << brightness << "% - Color: " << color << std::endl;
    }
};

class SmartThermostat : public SmartDevice {
private:
    float currentTemp;
    float targetTemp;
    std::string mode;
    
public:
    SmartThermostat(const std::string& id, const std::string& name,
                    const std::string& location, float currentTemp = 22.0f,
                    float targetTemp = 22.0f, const std::string& mode = "Auto")
        : SmartDevice(id, name, DeviceType::THERMOSTAT, location),
          currentTemp(currentTemp), targetTemp(targetTemp), mode(mode) {}
    
    void setTargetTemperature(float temp) {
        targetTemp = std::max(10.0f, std::min(35.0f, temp));
    }
    
    void updateCurrentTemperature(float temp) {
        currentTemp = temp;
    }
    
    void getStatusReport() const override {
        std::cout << name << " - Current: " << currentTemp << "°C - Target: " << targetTemp << "°C - Mode: " << mode << std::endl;
    }
};

class SmartDoorLock : public SmartDevice {
private:
    bool isLocked;
    std::string accessCode;
    
public:
    SmartDoorLock(const std::string& id, const std::string& name,
                  const std::string& location, const std::string& accessCode = "1234")
        : SmartDevice(id, name, DeviceType::DOOR_LOCK, location),
          isLocked(true), accessCode(accessCode) {}
    
    bool unlock(const std::string& code) {
        if (code == accessCode) {
            isLocked = false;
            return true;
        }
        return false;
    }
    
    void getStatusReport() const override {
        std::cout << name << " - Lock: " << (isLocked ? "Locked" : "Unlocked") << std::endl;
    }
};

class Scene {
private:
    std::string sceneName;
    std::vector<std::shared_ptr<SmartDevice>> devices;
    
public:
    Scene(const std::string& name) : sceneName(name) {}
    
    void addDevice(std::shared_ptr<SmartDevice> device) {
        devices.push_back(device);
    }
    
    void activate() {
        for (const auto& device : devices) {
            device->turnOn();
        }
    }
    
    void deactivate() {
        for (const auto& device : devices) {
            device->turnOff();
        }
    }
};

class SmartHomeSystem {
private:
    std::map<std::string, std::shared_ptr<SmartDevice>> devices;
    std::map<std::string, std::shared_ptr<Scene>> scenes;
    
public:
    void addDevice(std::shared_ptr<SmartDevice> device) {
        devices[device->getId()] = device;
    }
    
    void addScene(std::shared_ptr<Scene> scene) {
        scenes[scene->sceneName] = scene;
    }
    
    void activateScene(const std::string& sceneName) {
        if (scenes.find(sceneName) != scenes.end()) {
            scenes[sceneName]->activate();
        }
    }
    
    void generateReport() const {
        std::cout << "=== Report ===" << std::endl;
        std::cout << "Devices: " << devices.size() << std::endl;
        std::cout << "Scenes: " << scenes.size() << std::endl;
        
        for (const auto& pair : devices) {
            pair.second->getStatusReport();
        }
    }
};

int main() {
    SmartHomeSystem homeSystem;
    
    auto livingRoomLight = std::make_shared<SmartLight>("LGT001", "Living Room Light", "Living Room", 75);
    auto kitchenLight = std::make_shared<SmartLight>("LGT002", "Kitchen Light", "Kitchen", 100);
    auto mainThermostat = std::make_shared<SmartThermostat>("THR001", "Main Thermostat", "Hallway", 22.0f, 22.0f);
    auto frontDoorLock = std::make_shared<SmartDoorLock>("LCK001", "Front Door Lock", "Front Door", "1234");
    
    homeSystem.addDevice(livingRoomLight);
    homeSystem.addDevice(kitchenLight);
    homeSystem.addDevice(mainThermostat);
    homeSystem.addDevice(frontDoorLock);
    
    auto eveningScene = std::make_shared<Scene>("Evening");
    auto morningScene = std::make_shared<Scene>("Morning");
    
    eveningScene->addDevice(livingRoomLight);
    eveningScene->addDevice(kitchenLight);
    
    morningScene->addDevice(kitchenLight);
    morningScene->addDevice(mainThermostat);
    
    homeSystem.addScene(eveningScene);
    homeSystem.addScene(morningScene);
    
    livingRoomLight->turnOn();
    livingRoomLight->setBrightness(80);
    
    kitchenLight->turnOn();
    
    mainThermostat->turnOn();
    mainThermostat->setTargetTemperature(24.0f);
    
    homeSystem.activateScene("Morning");
    
    homeSystem.generateReport();
    
    return 0;
}