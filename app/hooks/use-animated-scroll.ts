import { useScrollContext } from '@/context/scroll-context';
import { useAnimatedScrollHandler } from 'react-native-reanimated';

export function useAnimatedScroll() {
  const { scrollY } = useScrollContext();

  const scrollHandler = useAnimatedScrollHandler({
    onScroll: (event) => {
      scrollY.value = event.contentOffset.y;
    },
  });

  return scrollHandler;
}
