import React, { useState, FC, ReactNode } from 'react';
import { View, Text, TouchableOpacity, StyleSheet, ScrollView, ViewStyle, TextStyle } from 'react-native';
import { MotiView } from 'moti';
import { Easing } from 'react-native-reanimated';

interface DropdownItem {
  label: string;
  value: string | number;
}

interface DropdownTriggerProps {
  label: string;
  isOpen: boolean;
  onPress: () => void;
  disabled?: boolean;
  style?: ViewStyle;
  textStyle?: TextStyle;
  children?: ReactNode;
}

interface DropdownMenuItemProps {
  item: DropdownItem;
  children?: ReactNode;
  isSelected: boolean;
  onPress: () => void;
  index: number;
  style?: ViewStyle;
  textStyle?: TextStyle;
}

interface DropdownMenuProps {
  isOpen: boolean;
  children: ReactNode;
  style?: ViewStyle;
}

interface DropdownContextType {
  isOpen: boolean;
  selectedValue: string | number | null;
  onSelect: (item: DropdownItem) => void;
}

const DropdownTrigger: FC<DropdownTriggerProps> = ({
  label,
  isOpen,
  onPress,
  disabled = false,
  style = {},
  textStyle = {},
  children,
}) => (
  <TouchableOpacity
    onPress={onPress}
    disabled={disabled}
    activeOpacity={0.7}
  >
    <View style={[styles.trigger, style, disabled && styles.disabled]}>
      <Text style={[styles.triggerText, textStyle]}>{label}</Text>
      <MotiView
        from={{ rotate: '0deg' }}
        animate={{ rotate: isOpen ? '180deg' : '0deg' }}
        transition={{
          type: 'timing',
          duration: 300,
          easing: Easing.inOut(Easing.ease),
        }}
      >
        {children}
      </MotiView>
    </View>
  </TouchableOpacity>
);

const DropdownMenuItem: FC<DropdownMenuItemProps> = ({
  item,
  isSelected,
  onPress,
  index,
  style = {},
  textStyle = {},
  children,
}) => (
  <MotiView
    from={{
      opacity: 0,
      translateY: -10,
    }}
    animate={{
      opacity: 1,
      translateY: 0,
    }}
    transition={{
      type: 'timing',
      duration: 150,
      delay: index * 30,
      easing: Easing.out(Easing.ease),
    }}
  >
    <TouchableOpacity
      onPress={onPress}
      activeOpacity={0.6}
    >
      <View
        style={[
          styles.menuItem,
          style,
          isSelected && styles.selectedMenuItem,
        ]}
      >
        <Text
          style={[
            styles.menuItemText,
            textStyle,
            isSelected && styles.selectedMenuItemText,
          ]}
        >
          {item.label}
        </Text>
      </View>
    </TouchableOpacity>
  </MotiView>
);

// Dropdown Menu Component
const DropdownMenu: FC<DropdownMenuProps> = ({
  isOpen,
  children,
  style = {},
}) => {
  if (!isOpen) return null;

  return (
    <MotiView
      from={{
        opacity: 0,
        scale: 0.9,
      }}
      animate={{
        opacity: 1,
        scale: 1,
      }}
      exit={{
        opacity: 0,
        scale: 0.9,
      }}
      transition={{
        type: 'timing',
        duration: 200,
        easing: Easing.inOut(Easing.ease),
      }}
      style={[styles.menu, style]}
    >
      <ScrollView
        scrollEnabled={true}
        nestedScrollEnabled={true}
        style={styles.menuContent}
      >
        {children}
      </ScrollView>
    </MotiView>
  );
};

interface DropdownProps {
  items: DropdownItem[];
  onSelect: (item: DropdownItem) => void;
  placeholder?: string;
  selectedValue?: string | number | null;
  disabled?: boolean;
  triggerStyle?: ViewStyle;
  triggerTextStyle?: TextStyle;
  menuStyle?: ViewStyle;
  menuItemStyle?: ViewStyle;
  menuItemTextStyle?: TextStyle;
}

const Dropdown: FC<DropdownProps> = ({
  items,
  onSelect,
  placeholder = 'Select an option',
  selectedValue = null,
  disabled = false,
  triggerStyle = {},
  triggerTextStyle = {},
  menuStyle = {},
  menuItemStyle = {},
  menuItemTextStyle = {},
}) => {
  const [isOpen, setIsOpen] = useState<boolean>(false);

  const handleSelectItem = (item: DropdownItem): void => {
    onSelect(item);
    setIsOpen(false);
  };

  const selectedLabel: string = 
    items.find(item => item.value === selectedValue)?.label || placeholder;

  return (
    <View style={styles.container}>
      <DropdownTrigger
        label={selectedLabel}
        isOpen={isOpen}
        onPress={() => !disabled && setIsOpen(!isOpen)}
        disabled={disabled}
        style={triggerStyle}
        textStyle={triggerTextStyle}
      />

      <DropdownMenu
        isOpen={isOpen}
        style={menuStyle}
      >
        {items.map((item, index) => (
          <DropdownMenuItem
            key={item.value}
            item={item}
            isSelected={selectedValue === item.value}
            onPress={() => handleSelectItem(item)}
            index={index}
            style={menuItemStyle}
            textStyle={menuItemTextStyle}
          />
        ))}
      </DropdownMenu>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    width: '100%',
    zIndex: 1000,
  },
  trigger: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 16,
    paddingVertical: 12,
    backgroundColor: '#fff',
    borderRadius: 8,
    borderWidth: 1,
    borderColor: '#ddd',
  },
  triggerText: {
    fontSize: 16,
    color: '#333',
    flex: 1,
  },
  chevron: {
    fontSize: 12,
    color: '#666',
    marginLeft: 8,
  },
  disabled: {
    opacity: 0.5,
  },
  menu: {
    backgroundColor: '#fff',
    borderRadius: 8,
    borderWidth: 1,
    borderColor: '#ddd',
    marginTop: 4,
    overflow: 'hidden',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 5,
  },
  menuContent: {
    maxHeight: 250,
  },
  menuItem: {
    paddingHorizontal: 16,
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: '#f0f0f0',
  },
  menuItemText: {
    fontSize: 14,
    color: '#333',
  },
  selectedMenuItem: {
    backgroundColor: '#f5f5f5',
  },
  selectedMenuItemText: {
    fontWeight: '600',
    color: '#007AFF',
  },
});

export default Dropdown;
export { DropdownTrigger, DropdownMenu, DropdownMenuItem };
export type { DropdownProps, DropdownTriggerProps, DropdownMenuProps, DropdownMenuItemProps, DropdownItem };