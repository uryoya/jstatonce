import java.io.*;

public class HeapTest {
    public static void main(String[] args) throws InterruptedException {
        int i = 0;
        while (i < 100) {
            Thread.sleep(100);
            /* 3MBの短命オブジェクトを作り続ける */
            StringBuffer tempStr = new StringBuffer(3000000);
            i++;
        }
    }
}
