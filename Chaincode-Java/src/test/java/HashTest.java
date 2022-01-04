import org.junit.jupiter.api.Test;

import java.time.LocalDate;

import static org.assertj.core.api.Assertions.assertThat;
public final class HashTest {

    @Test
    public void isEqual(){
        Hash hash = new Hash("hashashhashahs", LocalDate.now().toString(), "pippo");
        assertThat(hash).isEqualTo(hash);
    }

    @Test
    public void nonEqual(){
        Hash hash = new Hash("str1", LocalDate.parse("2021-01-01").toString(), "pippo");
        Hash hash2 = new Hash("str2", LocalDate.now().toString(), "paperino");
        assertThat(hash).isNotEqualTo(hash2);
    }
}