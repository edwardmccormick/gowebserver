import FontAwesome from "@expo/vector-icons/FontAwesome";
import MaterialIcons from "@expo/vector-icons/MaterialIcons";
import { StatusBar } from 'expo-status-bar';
import { StyleSheet, Text, View, Image } from 'react-native';

export default function App() {
  return (
    <View style={styles.container}>
      <Text>Back to the fucking beginning! lol....hello world</Text>
      <Text>Welcome to React Native with Expo!</Text>
      <Text>Let's build something amazing together.</Text>
      <Text>Enjoy coding!</Text>
      <Text>Keep pushing forward!</Text>
      <Text>Stay curious and keep learning!</Text>
      <Text>Remember, every expert was once a beginner.</Text>
      <Text>Embrace the journey of growth!</Text>
      <Text>Happy coding!</Text>
      <Text>Let's make the most of this!</Text>
      <Text>Cheers to new adventures!</Text>
      <Text>Keep smiling and coding!</Text>
      <FontAwesome name="apple" size={25} />
      <MaterialIcons name="star" color="blue" size={25} />
      {/* Create a button */}
      <FontAwesome.Button
        name="facebook"
        backgroundColor="#3b5998"
        onPress={() => {}}
      >
        Login with Facebook
      </FontAwesome.Button>
      <StatusBar style="auto" />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    alignItems: 'center',
    justifyContent: 'center',
  },
});
