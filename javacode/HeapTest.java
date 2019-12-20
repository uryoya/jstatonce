import java.io.*;

public class HeapTest {
    public static void main(String[] args) throws InterruptedException {
        int i = 0;
        while (i < 100) {
            Thread.sleep(100);
            /* 3MBの短命オブジェクトを作り続ける */
            StringBuffer tempStr = new StringBuffer(3000000);
            System.out.printf("stdout: %d\n", i);
            System.err.printf("stderr: %d\n", i);
            i++;
        }
    }
}
