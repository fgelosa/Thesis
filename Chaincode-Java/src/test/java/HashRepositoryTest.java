import org.hyperledger.fabric.contract.Context;
import org.hyperledger.fabric.shim.ChaincodeStub;
import org.hyperledger.fabric.shim.ChaincodeException;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.catchThrowable;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;


public final class HashRepositoryTest {

    @Nested
    class InvokeQueryAgreementTransaction {

        @Test
        public void whenAgreementExists() {
            HashRepository contract = new HashRepository();
            Context ctx = mock(Context.class);
            ChaincodeStub stub = mock(ChaincodeStub.class);
            when(ctx.getStub()).thenReturn(stub);
            when(stub.getStringState("hash000"))
                    .thenReturn("{\"hash\":\"hashashashash\",\"date\":\"2021-01-01\",\"user\":\"franco\"}");

            Hash hash = contract.getHash(ctx, "hash000");

            assertThat(hash.getHashString())
                    .isEqualTo("hashashashash");
            assertThat(hash.getDate())
                    .isEqualTo("2021-01-01");
            assertThat(hash.getUser())
                    .isEqualTo("franco");
        }

        @Test
        public void whenCarDoesNotExist() {
            HashRepository contract = new HashRepository();
            Context ctx = mock(Context.class);
            ChaincodeStub stub = mock(ChaincodeStub.class);
            when(ctx.getStub()).thenReturn(stub);
            when(stub.getStringState("hash000")).thenReturn("");

            Throwable thrown = catchThrowable(() -> {
                contract.getHash(ctx, "hash000");
            });

            assertThat(thrown).isInstanceOf(ChaincodeException.class).hasNoCause()
                    .hasMessage("Hash record hash000 does not exist");
        }

    }

}
